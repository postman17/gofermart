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

func TestRegisterUser_Success(t *testing.T) {
	repos := &mockDBRepository{
		registerUserFunc: func(ctx context.Context, login, password string) error {
			return nil
		},
	}
	handler := RegisterUser(context.Background(), repos)

	body, _ := json.Marshal(models.AuthUser{Login: "testuser", Password: "testpass"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(string(body)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
}

func TestRegisterUser_WrongMethod(t *testing.T) {
	repos := &mockDBRepository{}
	handler := RegisterUser(context.Background(), repos)

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		req := httptest.NewRequest(method, "/api/user/register", nil)
		rw := httptest.NewRecorder()
		handler(rw, req)

		if rw.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected status 405, got %d", method, rw.Code)
		}
	}
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	repos := &mockDBRepository{}
	handler := RegisterUser(context.Background(), repos)

	req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader("not json"))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rw.Code)
	}
}

func TestRegisterUser_EmptyLoginOrPassword(t *testing.T) {
	repos := &mockDBRepository{}
	handler := RegisterUser(context.Background(), repos)

	tests := []struct {
		name     string
		body     models.AuthUser
		expected int
	}{
		{"empty login", models.AuthUser{Login: "", Password: "pass"}, http.StatusBadRequest},
		{"empty password", models.AuthUser{Login: "user", Password: ""}, http.StatusBadRequest},
		{"both empty", models.AuthUser{Login: "", Password: ""}, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(string(body)))
			rw := httptest.NewRecorder()
			handler(rw, req)

			if rw.Code != tt.expected {
				t.Errorf("expected status %d, got %d", tt.expected, rw.Code)
			}
		})
	}
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	repos := &mockDBRepository{
		registerUserFunc: func(ctx context.Context, login, password string) error {
			return errorInternal.ErrUserAlreadyExists
		},
	}
	handler := RegisterUser(context.Background(), repos)

	body, _ := json.Marshal(models.AuthUser{Login: "existing", Password: "pass"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(string(body)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", rw.Code)
	}
}

func TestRegisterUser_RepositoryError(t *testing.T) {
	repos := &mockDBRepository{
		registerUserFunc: func(ctx context.Context, login, password string) error {
			return errors.New("db error")
		},
	}
	handler := RegisterUser(context.Background(), repos)

	body, _ := json.Marshal(models.AuthUser{Login: "user", Password: "pass"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(string(body)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}
