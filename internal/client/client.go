package client

import (
	"context"

	models "github.com/postman17/gofermart/internal/model"
)

type AccrualSystemClientRepository interface {
	GetAccrual(orderId string) (models.OrderResult, error)
}

type AccrualSystemClient struct {
	Url string
}

func NewAccrualSystemClient(ctx context.Context, url string) AccrualSystemClientRepository {
	return &AccrualSystemClient{
		Url: url,
	}
}
