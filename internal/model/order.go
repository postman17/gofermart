package model

import "time"

type OrderCheckResult struct {
	Exists     bool // Существует ли заказ в базе
	IsSameUser bool // Принадлежит ли он текущему пользователю
}

type OrderListItem struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

type OrderResult struct {
	Status  string
	Accrual float64
}
