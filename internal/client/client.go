package client

import (
	"context"
	"net/http"
	"time"

	models "github.com/postman17/gofermart/internal/model"
)

type AccrualSystemClientRepository interface {
	GetAccrual(orderId string) (models.OrderResult, error)
}

type AccrualSystemClient struct {
	Url string
	c   http.Client
}

func NewAccrualSystemClient(ctx context.Context, url string) AccrualSystemClientRepository {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	return &AccrualSystemClient{
		Url: url,
		c:   client,
	}
}
