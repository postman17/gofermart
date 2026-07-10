package handlers

import (
	"context"

	models "github.com/postman17/gofermart/internal/model"
)

type mockDBRepository struct {
	registerUserFunc          func(ctx context.Context, login, password string) error
	userExistsFunc            func(ctx context.Context, login string) (bool, error)
	authenticateUserFunc      func(ctx context.Context, login, password string) (int64, error)
	createOrUpdateSessionFunc func(ctx context.Context, userID int64) (string, error)
	getUserIDByTokenFunc      func(ctx context.Context, token string) (int64, error)
	checkOrderOwnershipFunc   func(ctx context.Context, orderNumber string, currentUserID int64) (models.OrderCheckResult, error)
	createOrderFunc           func(ctx context.Context, userId int64, number string, result models.OrderResult) (int, error)
	getOrdersByUserIDFunc     func(ctx context.Context, userID int64) ([]models.OrderListItem, error)
	getBalanceFunc            func(ctx context.Context, userID int64) (models.Balance, error)
	accrueBalanceFunc         func(ctx context.Context, userId int64, amount float64) error
	withdrawFunc              func(ctx context.Context, orderID string, userID int64, amount int64) error
	getUserWithdrawalsFunc    func(ctx context.Context, userID int64) ([]models.Withdraw, error)
}

func (m *mockDBRepository) RegisterUser(ctx context.Context, login, password string) error {
	if m.registerUserFunc != nil {
		return m.registerUserFunc(ctx, login, password)
	}
	return nil
}

func (m *mockDBRepository) UserExists(ctx context.Context, login string) (bool, error) {
	if m.userExistsFunc != nil {
		return m.userExistsFunc(ctx, login)
	}
	return false, nil
}

func (m *mockDBRepository) AuthenticateUser(ctx context.Context, login, password string) (int64, error) {
	if m.authenticateUserFunc != nil {
		return m.authenticateUserFunc(ctx, login, password)
	}
	return 0, nil
}

func (m *mockDBRepository) CreateOrUpdateSession(ctx context.Context, userID int64) (string, error) {
	if m.createOrUpdateSessionFunc != nil {
		return m.createOrUpdateSessionFunc(ctx, userID)
	}
	return "test-token", nil
}

func (m *mockDBRepository) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	if m.getUserIDByTokenFunc != nil {
		return m.getUserIDByTokenFunc(ctx, token)
	}
	return 1, nil
}

func (m *mockDBRepository) CheckOrderOwnership(ctx context.Context, orderNumber string, currentUserID int64) (models.OrderCheckResult, error) {
	if m.checkOrderOwnershipFunc != nil {
		return m.checkOrderOwnershipFunc(ctx, orderNumber, currentUserID)
	}
	return models.OrderCheckResult{}, nil
}

func (m *mockDBRepository) CreateOrder(ctx context.Context, userId int64, number string, result models.OrderResult) (int, error) {
	if m.createOrderFunc != nil {
		return m.createOrderFunc(ctx, userId, number, result)
	}
	return 0, nil
}

func (m *mockDBRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.OrderListItem, error) {
	if m.getOrdersByUserIDFunc != nil {
		return m.getOrdersByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockDBRepository) GetBalance(ctx context.Context, userID int64) (models.Balance, error) {
	if m.getBalanceFunc != nil {
		return m.getBalanceFunc(ctx, userID)
	}
	return models.Balance{}, nil
}

func (m *mockDBRepository) AccrueBalance(ctx context.Context, userId int64, amount float64) error {
	if m.accrueBalanceFunc != nil {
		return m.accrueBalanceFunc(ctx, userId, amount)
	}
	return nil
}

func (m *mockDBRepository) Withdraw(ctx context.Context, orderID string, userID int64, amount int64) error {
	if m.withdrawFunc != nil {
		return m.withdrawFunc(ctx, orderID, userID, amount)
	}
	return nil
}

func (m *mockDBRepository) GetUserWithdrawals(ctx context.Context, userID int64) ([]models.Withdraw, error) {
	if m.getUserWithdrawalsFunc != nil {
		return m.getUserWithdrawalsFunc(ctx, userID)
	}
	return nil, nil
}
