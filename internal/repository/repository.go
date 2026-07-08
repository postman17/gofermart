package repository

import (
	"context"
	"database/sql"

	models "github.com/postman17/gofermart/internal/model"
)

type DBRepository interface {
	RegisterUser(login string, password string) error
	UserExists(login string) (bool, error)
	AuthenticateUser(login string, password string) (int64, error)
	CreateOrUpdateSession(userID int64) (string, error)
	GetUserIDByToken(token string) (int64, error)
	CheckOrderOwnership(orderNumber string, currentUserID int64) (models.OrderCheckResult, error)
	CreateOrder(userId int64, number string, result models.OrderResult) (int, error)
	GetOrdersByUserID(userID int64) ([]models.OrderListItem, error)
	GetBalance(userID int64) (models.Balance, error)
	AccrueBalance(userId int64, amount int64) error
	Withdraw(orderID string, userID int64, amount int64) error
	GetUserWithdrawals(userID int64) ([]models.Withdraw, error)
}

type DBStorage struct {
	db  *sql.DB
	ctx context.Context
}

func NewDBStorage(ctx context.Context, db *sql.DB) DBRepository {
	return &DBStorage{
		db:  db,
		ctx: ctx,
	}
}
