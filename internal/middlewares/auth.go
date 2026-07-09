package middlewares

import (
	"context"
	"net/http"
	"strings"

	repo "github.com/postman17/gofermart/internal/repository"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(ctx context.Context, repos repo.DBRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token := parts[1]

			userID, err := repos.GetUserIDByToken(ctx, token)
			if err != nil {
				http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
