package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	error "github.com/postman17/gofermart/internal/errors"
	models "github.com/postman17/gofermart/internal/model"
	repo "github.com/postman17/gofermart/internal/repository"
)

func RegisterUser(ctx context.Context, repos repo.DBRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var req models.AuthUser
		if err := json.Unmarshal(body, &req); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Login == "" || req.Password == "" {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		err = repos.RegisterUser(ctx, req.Login, req.Password)
		if err != nil {
			if errors.Is(err, error.ErrUserAlreadyExists) {
				rw.WriteHeader(http.StatusConflict)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("{}"))
	}
}
