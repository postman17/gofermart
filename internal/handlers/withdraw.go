package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	errorInternal "github.com/postman17/gofermart/internal/errors"
	models "github.com/postman17/gofermart/internal/model"
	repo "github.com/postman17/gofermart/internal/repository"
)

func Withdraw(repos repo.DBRepository) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
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

		var req models.Withdraw
		if err := json.Unmarshal(body, &req); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Order == "" || req.TotalSum == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		if !IsValidLuhn(req.Order) {
			rw.WriteHeader(http.StatusUnprocessableEntity)
			return

		}
		userID, ok := r.Context().Value("userID").(int64)

		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = repos.Withdraw(req.Order, userID, req.TotalSum)
		if err != nil {
			if errors.Is(err, errorInternal.ErrInsufficientFunds) {
				rw.WriteHeader(http.StatusPaymentRequired)
			}
			if errors.Is(err, errorInternal.ErrInvalidAmount) {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}

func UserWithdrawals(repos repo.DBRepository) http.HandlerFunc {
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

		resp, err := repos.GetUserWithdrawals(userID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(resp) == 0 {
			rw.WriteHeader(http.StatusNoContent)
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
