package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"unicode"

	clients "github.com/postman17/gofermart/internal/client"
	repo "github.com/postman17/gofermart/internal/repository"
)

func IsValidLuhn(s string) bool {
	var sum int
	parity := len(s) % 2

	// Идем по строке справа налево
	for i := len(s) - 1; i >= 0; i-- {
		r := rune(s[i])

		// Игнорируем пробелы, если они есть
		if unicode.IsSpace(r) {
			// Корректируем четность из-за пропущенного символа
			if i%2 == parity {
				parity = 1 - parity
			}
			continue
		}

		// Проверяем, что символ является цифрой
		if !unicode.IsDigit(r) {
			return false
		}

		// Конвертируем символ в число
		digit := int(r - '0')

		// Удваиваем каждую вторую цифру
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
	}

	// Номер корректен, если сумма делится на 10 без остатка
	return sum%10 == 0 && sum > 0
}

func AddOrder(repos repo.DBRepository, client clients.AccrualSystemClient) http.HandlerFunc {
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

		receivedText := string(body)
		if !IsValidLuhn(receivedText) {
			rw.WriteHeader(http.StatusBadRequest)
			return

		}
		userID, ok := r.Context().Value("userID").(int64)

		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		checkResult, err := repos.CheckOrderOwnership(receivedText, userID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if checkResult.Exists {
			if checkResult.IsSameUser {
				rw.WriteHeader(http.StatusOK)
				return
			} else {
				rw.WriteHeader(http.StatusConflict)
				return
			}
		}

		result, err := client.GetAccrual(receivedText)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		repos.CreateOrder(userID, receivedText, result)
		repos.AccrueBalance(userID, int64(result.Accrual))
		rw.WriteHeader(http.StatusAccepted)
	}
}

func GetOrders(repos repo.DBRepository) http.HandlerFunc {
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

		resp, err := repos.GetOrdersByUserID(userID)
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
