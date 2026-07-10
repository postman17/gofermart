package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	clients "github.com/postman17/gofermart/internal/client"
	models "github.com/postman17/gofermart/internal/model"
)

func TestIsValidLuhn(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"valid number 79927398713", "79927398713", true},
		{"valid number 49927398716", "49927398716", true},
		{"valid single digit not valid", "0", false},
		{"invalid number 79927398714", "79927398714", false},
		{"invalid number 49927398717", "49927398717", false},
		{"empty string", "", false},
		{"letters", "abc", false},
		{"mixed", "7992abc98713", false},
		{"valid 4561261212345467", "4561261212345467", true},
		{"invalid 4561261212345468", "4561261212345468", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidLuhn(tt.input)
			if result != tt.valid {
				t.Errorf("IsValidLuhn(%q) = %v, want %v", tt.input, result, tt.valid)
			}
		})
	}
}

func newAccrualTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestAddOrder_Success(t *testing.T) {
	accrualServer := newAccrualTestServer(func(w http.ResponseWriter, r *http.Request) {
		resp := models.OrderResponse{
			Order:  "79927398713",
			Status: "PROCESSED",
		}
		accrual := 100.0
		resp.Accrual = &accrual
		json.NewEncoder(w).Encode(resp)
	})
	defer accrualServer.Close()

	repos := &mockDBRepository{
		checkOrderOwnershipFunc: func(ctx context.Context, orderNumber string, currentUserID int64) (models.OrderCheckResult, error) {
			return models.OrderCheckResult{Exists: false}, nil
		},
		createOrderFunc: func(ctx context.Context, userId int64, number string, result models.OrderResult) (int, error) {
			return 0, nil
		},
		accrueBalanceFunc: func(ctx context.Context, userId int64, amount float64) error {
			return nil
		},
	}

	client := &clients.AccrualSystemClient{Url: accrualServer.URL}
	handler := AddOrder(context.Background(), repos, client)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("79927398713"))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusAccepted {
		t.Errorf("expected status 202, got %d", rw.Code)
	}
}

func TestAddOrder_WrongMethod(t *testing.T) {
	repos := &mockDBRepository{}
	client := &clients.AccrualSystemClient{}
	handler := AddOrder(context.Background(), repos, client)

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/api/user/orders", nil)
		rw := httptest.NewRecorder()
		handler(rw, req)

		if rw.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected status 405, got %d", method, rw.Code)
		}
	}
}

func TestAddOrder_InvalidLuhn(t *testing.T) {
	repos := &mockDBRepository{}
	client := &clients.AccrualSystemClient{}
	handler := AddOrder(context.Background(), repos, client)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("12345678901"))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rw.Code)
	}
}

func TestAddOrder_MissingUserID(t *testing.T) {
	repos := &mockDBRepository{}
	client := &clients.AccrualSystemClient{}
	handler := AddOrder(context.Background(), repos, client)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("79927398713"))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestAddOrder_OrderAlreadyExistsSameUser(t *testing.T) {
	repos := &mockDBRepository{
		checkOrderOwnershipFunc: func(ctx context.Context, orderNumber string, currentUserID int64) (models.OrderCheckResult, error) {
			return models.OrderCheckResult{Exists: true, IsSameUser: true}, nil
		},
	}
	client := &clients.AccrualSystemClient{}
	handler := AddOrder(context.Background(), repos, client)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("79927398713"))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
}

func TestAddOrder_OrderAlreadyExistsDifferentUser(t *testing.T) {
	repos := &mockDBRepository{
		checkOrderOwnershipFunc: func(ctx context.Context, orderNumber string, currentUserID int64) (models.OrderCheckResult, error) {
			return models.OrderCheckResult{Exists: true, IsSameUser: false}, nil
		},
	}
	client := &clients.AccrualSystemClient{}
	handler := AddOrder(context.Background(), repos, client)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("79927398713"))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", rw.Code)
	}
}

func TestAddOrder_CheckOrderOwnershipError(t *testing.T) {
	repos := &mockDBRepository{
		checkOrderOwnershipFunc: func(ctx context.Context, orderNumber string, currentUserID int64) (models.OrderCheckResult, error) {
			return models.OrderCheckResult{}, errors.New("db error")
		},
	}
	client := &clients.AccrualSystemClient{}
	handler := AddOrder(context.Background(), repos, client)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("79927398713"))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestAddOrder_AccrualServerError(t *testing.T) {
	repos := &mockDBRepository{
		checkOrderOwnershipFunc: func(ctx context.Context, orderNumber string, currentUserID int64) (models.OrderCheckResult, error) {
			return models.OrderCheckResult{Exists: false}, nil
		},
	}
	client := &clients.AccrualSystemClient{Url: "http://localhost:0"}
	handler := AddOrder(context.Background(), repos, client)

	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", strings.NewReader("79927398713"))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestGetOrders_Success(t *testing.T) {
	orders := []models.OrderListItem{
		{ID: 1, Number: "79927398713", Status: "PROCESSED"},
		{ID: 2, Number: "49927398716", Status: "NEW"},
	}
	repos := &mockDBRepository{
		getOrdersByUserIDFunc: func(ctx context.Context, userID int64) ([]models.OrderListItem, error) {
			return orders, nil
		},
	}
	handler := GetOrders(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
}

func TestGetOrders_WrongMethod(t *testing.T) {
	repos := &mockDBRepository{}
	handler := GetOrders(context.Background(), repos)

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/api/user/orders", nil)
		rw := httptest.NewRecorder()
		handler(rw, req)

		if rw.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected status 405, got %d", method, rw.Code)
		}
	}
}

func TestGetOrders_MissingUserID(t *testing.T) {
	repos := &mockDBRepository{}
	handler := GetOrders(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestGetOrders_RepositoryError(t *testing.T) {
	repos := &mockDBRepository{
		getOrdersByUserIDFunc: func(ctx context.Context, userID int64) ([]models.OrderListItem, error) {
			return nil, errors.New("db error")
		},
	}
	handler := GetOrders(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}
