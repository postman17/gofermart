package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	errorInternal "github.com/postman17/gofermart/internal/errors"
	models "github.com/postman17/gofermart/internal/model"
)

func TestWithdraw_Success(t *testing.T) {
	repos := &mockDBRepository{
		withdrawFunc: func(ctx context.Context, orderID string, userID int64, amount int64) error {
			return nil
		},
	}
	handler := Withdraw(context.Background(), repos)

	body, _ := json.Marshal(models.Withdraw{Order: "79927398713", TotalSum: 100})
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(string(body)))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
}

func TestWithdraw_WrongMethod(t *testing.T) {
	repos := &mockDBRepository{}
	handler := Withdraw(context.Background(), repos)

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/api/user/balance/withdraw", nil)
		rw := httptest.NewRecorder()
		handler(rw, req)

		if rw.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected status 405, got %d", method, rw.Code)
		}
	}
}

func TestWithdraw_InvalidJSON(t *testing.T) {
	repos := &mockDBRepository{}
	handler := Withdraw(context.Background(), repos)

	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader("not json"))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rw.Code)
	}
}

func TestWithdraw_EmptyOrderOrZeroSum(t *testing.T) {
	repos := &mockDBRepository{}
	handler := Withdraw(context.Background(), repos)

	tests := []struct {
		name string
		body models.Withdraw
	}{
		{"empty order", models.Withdraw{Order: "", TotalSum: 100}},
		{"zero sum", models.Withdraw{Order: "79927398713", TotalSum: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(string(body)))
			req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
			rw := httptest.NewRecorder()
			handler(rw, req)

			if rw.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", rw.Code)
			}
		})
	}
}

func TestWithdraw_InvalidLuhnOrder(t *testing.T) {
	repos := &mockDBRepository{}
	handler := Withdraw(context.Background(), repos)

	body, _ := json.Marshal(models.Withdraw{Order: "12345678901", TotalSum: 100})
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(string(body)))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", rw.Code)
	}
}

func TestWithdraw_MissingUserID(t *testing.T) {
	repos := &mockDBRepository{}
	handler := Withdraw(context.Background(), repos)

	body, _ := json.Marshal(models.Withdraw{Order: "79927398713", TotalSum: 100})
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(string(body)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	repos := &mockDBRepository{
		withdrawFunc: func(ctx context.Context, orderID string, userID int64, amount int64) error {
			return errorInternal.ErrInsufficientFunds
		},
	}
	handler := Withdraw(context.Background(), repos)

	body, _ := json.Marshal(models.Withdraw{Order: "79927398713", TotalSum: 9999})
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(string(body)))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusPaymentRequired {
		t.Errorf("expected status 402, got %d", rw.Code)
	}
}

func TestWithdraw_InvalidAmount(t *testing.T) {
	repos := &mockDBRepository{
		withdrawFunc: func(ctx context.Context, orderID string, userID int64, amount int64) error {
			return errorInternal.ErrInvalidAmount
		},
	}
	handler := Withdraw(context.Background(), repos)

	body, _ := json.Marshal(models.Withdraw{Order: "79927398713", TotalSum: 100})
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(string(body)))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rw.Code)
	}
}

func TestWithdraw_RepositoryError(t *testing.T) {
	repos := &mockDBRepository{
		withdrawFunc: func(ctx context.Context, orderID string, userID int64, amount int64) error {
			return errors.New("db error")
		},
	}
	handler := Withdraw(context.Background(), repos)

	body, _ := json.Marshal(models.Withdraw{Order: "79927398713", TotalSum: 100})
	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(string(body)))
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestUserWithdrawals_Success(t *testing.T) {
	withdrawals := []models.Withdraw{
		{Order: "79927398713", TotalSum: 100},
		{Order: "49927398716", TotalSum: 200},
	}
	repos := &mockDBRepository{
		getUserWithdrawalsFunc: func(ctx context.Context, userID int64) ([]models.Withdraw, error) {
			return withdrawals, nil
		},
	}
	handler := UserWithdrawals(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
}

func TestUserWithdrawals_WrongMethod(t *testing.T) {
	repos := &mockDBRepository{}
	handler := UserWithdrawals(context.Background(), repos)

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/api/user/withdrawals", nil)
		rw := httptest.NewRecorder()
		handler(rw, req)

		if rw.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected status 405, got %d", method, rw.Code)
		}
	}
}

func TestUserWithdrawals_MissingUserID(t *testing.T) {
	repos := &mockDBRepository{}
	handler := UserWithdrawals(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestUserWithdrawals_NoWithdrawals(t *testing.T) {
	repos := &mockDBRepository{
		getUserWithdrawalsFunc: func(ctx context.Context, userID int64) ([]models.Withdraw, error) {
			return []models.Withdraw{}, nil
		},
	}
	handler := UserWithdrawals(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", rw.Code)
	}
}

func TestUserWithdrawals_RepositoryError(t *testing.T) {
	repos := &mockDBRepository{
		getUserWithdrawalsFunc: func(ctx context.Context, userID int64) ([]models.Withdraw, error) {
			return nil, errors.New("db error")
		},
	}
	handler := UserWithdrawals(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}
