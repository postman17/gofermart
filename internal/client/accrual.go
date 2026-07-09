package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/postman17/gofermart/internal/model"
)

func (a *AccrualSystemClient) GetAccrual(orderId string) (models.OrderResult, error) {
	url := fmt.Sprintf("%s/api/orders/%s", a.Url, orderId)
	resp, err := a.c.Get(url)
	if err != nil {
		return models.OrderResult{}, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.OrderResult{}, fmt.Errorf("wrong response status: %d", resp.StatusCode)
	}

	var orderResp models.OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return models.OrderResult{}, fmt.Errorf("error unmarshal JSON: %w", err)
	}

	// Проверяем наличие поля accrual. Если его нет (nil), возвращаем 0.
	var accrualValue float64
	if orderResp.Accrual != nil {
		accrualValue = *orderResp.Accrual
	}

	return models.OrderResult{
		Status:  orderResp.Status,
		Accrual: accrualValue,
	}, nil
}
