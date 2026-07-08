package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/postman17/gofermart/internal/model"
)

func TestLoginUser_Success(t *testing.T) {
	repos := &mockDBRepository{
		authenticateUserFunc: func(login, password string) (int64, error) {
			return 1, nil
		},
		createOrUpdateSessionFunc: func(userID int64) (string, error) {
			return "test-jwt-token", nil
		},
	}
	handler := LoginUser(repos)

	body, _ := json.Marshal(models.AuthUser{Login: "testuser", Password: "testpass"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(string(body)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rw.Code)
	}
	authHeader := rw.Header().Get("Authorization")
	if authHeader != "Bearer test-jwt-token" {
		t.Errorf("expected Authorization header 'Bearer test-jwt-token', got '%s'", authHeader)
	}
}

func TestLoginUser_WrongMethod(t *testing.T) {
	repos := &mockDBRepository{}
	handler := LoginUser(repos)

	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/api/user/login", nil)
		rw := httptest.NewRecorder()
		handler(rw, req)

		if rw.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected status 405, got %d", method, rw.Code)
		}
	}
}

func TestLoginUser_InvalidJSON(t *testing.T) {
	repos := &mockDBRepository{}
	handler := LoginUser(repos)

	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader("not json"))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rw.Code)
	}
}

func TestLoginUser_EmptyLoginOrPassword(t *testing.T) {
	repos := &mockDBRepository{}
	handler := LoginUser(repos)

	tests := []struct {
		name string
		body models.AuthUser
	}{
		{"empty login", models.AuthUser{Login: "", Password: "pass"}},
		{"empty password", models.AuthUser{Login: "user", Password: ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(string(body)))
			rw := httptest.NewRecorder()
			handler(rw, req)

			if rw.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", rw.Code)
			}
		})
	}
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	repos := &mockDBRepository{
		authenticateUserFunc: func(login, password string) (int64, error) {
			return 0, nil
		},
	}
	handler := LoginUser(repos)

	body, _ := json.Marshal(models.AuthUser{Login: "testuser", Password: "wrongpass"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(string(body)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rw.Code)
	}
}

func TestLoginUser_AuthenticateError(t *testing.T) {
	repos := &mockDBRepository{
		authenticateUserFunc: func(login, password string) (int64, error) {
			return 0, errors.New("db error")
		},
	}
	handler := LoginUser(repos)

	body, _ := json.Marshal(models.AuthUser{Login: "testuser", Password: "testpass"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(string(body)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}

func TestLoginUser_SessionError(t *testing.T) {
	repos := &mockDBRepository{
		authenticateUserFunc: func(login, password string) (int64, error) {
			return 1, nil
		},
		createOrUpdateSessionFunc: func(userID int64) (string, error) {
			return "", errors.New("session error")
		},
	}
	handler := LoginUser(repos)

	body, _ := json.Marshal(models.AuthUser{Login: "testuser", Password: "testpass"})
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(string(body)))
	rw := httptest.NewRecorder()

	handler(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rw.Code)
	}
}
