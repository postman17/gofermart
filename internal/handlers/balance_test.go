package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/postman17/gofermart/internal/model"
)

func TestGetBalance_Success(t *testing.T) {
	expected := models.Balance{ID: 1, Current: 500, Withdrawn: 100}
	repos := &mockDBRepository{
		getBalanceFunc: func(ctx context.Context, userID int64) (models.Balance, error) {
			return expected, nil
		},
	}
	handler := GetBalance(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
}

func TestGetBalance_WrongMethod(t *testing.T) {
	repos := &mockDBRepository{}
	handler := GetBalance(context.Background(), repos)

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/api/user/balance", nil)
		rw := httptest.NewRecorder()
		handler(rw, req)

		if rw.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected status 405, got %d", method, rw.Code)
		}
	}
}

func TestGetBalance_MissingUserID(t *testing.T) {
	repos := &mockDBRepository{}
	handler := GetBalance(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestGetBalance_WrongUserIDType(t *testing.T) {
	repos := &mockDBRepository{}
	handler := GetBalance(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", "not-an-int64"))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestGetBalance_RepositoryError(t *testing.T) {
	repos := &mockDBRepository{
		getBalanceFunc: func(ctx context.Context, userID int64) (models.Balance, error) {
			return models.Balance{}, errors.New("db error")
		},
	}
	handler := GetBalance(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestGetBalance_ResponseJSON(t *testing.T) {
	expected := models.Balance{ID: 1, Current: 500, Withdrawn: 100}
	repos := &mockDBRepository{
		getBalanceFunc: func(ctx context.Context, userID int64) (models.Balance, error) {
			return expected, nil
		},
	}
	handler := GetBalance(context.Background(), repos)

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	req = req.WithContext(context.WithValue(req.Context(), "userID", int64(1)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	var result models.Balance
	if err := json.NewDecoder(rw.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result.Current != expected.Current {
		t.Errorf("expected current %d, got %d", expected.Current, result.Current)
	}
	if result.Withdrawn != expected.Withdrawn {
		t.Errorf("expected withdrawn %d, got %d", expected.Withdrawn, result.Withdrawn)
	}
}
