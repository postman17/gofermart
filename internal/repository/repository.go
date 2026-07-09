package repository

import (
	"context"
	"database/sql"

	models "github.com/postman17/gofermart/internal/model"
)

type DBRepository interface {
	RegisterUser(ctx context.Context, login string, password string) error
	UserExists(ctx context.Context, login string) (bool, error)
	AuthenticateUser(ctx context.Context, login string, password string) (int64, error)
	CreateOrUpdateSession(ctx context.Context, userID int64) (string, error)
	GetUserIDByToken(ctx context.Context, token string) (int64, error)
	CheckOrderOwnership(ctx context.Context, orderNumber string, currentUserID int64) (models.OrderCheckResult, error)
	CreateOrder(ctx context.Context, userId int64, number string, result models.OrderResult) (int, error)
	GetOrdersByUserID(ctx context.Context, userID int64) ([]models.OrderListItem, error)
	GetBalance(ctx context.Context, userID int64) (models.Balance, error)
	AccrueBalance(ctx context.Context, userId int64, amount int64) error
	Withdraw(ctx context.Context, orderID string, userID int64, amount int64) error
	GetUserWithdrawals(ctx context.Context, userID int64) ([]models.Withdraw, error)
}

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(ctx context.Context, db *sql.DB) DBRepository {
	return &DBStorage{
		db: db,
	}
}
