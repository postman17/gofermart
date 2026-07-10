package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	repo "github.com/postman17/gofermart/internal/repository"
)

func GetBalance(ctx context.Context, repos repo.DBRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		userID, ok := r.Context().Value("userID").(int64)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := repos.GetBalance(ctx, userID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		enc := json.NewEncoder(rw)
		if err := enc.Encode(resp); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}
