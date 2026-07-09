package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	models "github.com/postman17/gofermart/internal/model"
	repo "github.com/postman17/gofermart/internal/repository"
)

func LoginUser(ctx context.Context, repos repo.DBRepository) http.HandlerFunc {
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

		userID, err := repos.AuthenticateUser(ctx, req.Login, req.Password)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		} else if userID == 0 {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, err := repos.CreateOrUpdateSession(ctx, userID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Authorization", "Bearer "+token)

		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("{}"))
	}
}
