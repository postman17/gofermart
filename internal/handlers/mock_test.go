package handlers

import (
	models "github.com/postman17/gofermart/internal/model"
)

type mockDBRepository struct {
	registerUserFunc          func(login, password string) error
	userExistsFunc            func(login string) (bool, error)
	authenticateUserFunc      func(login, password string) (int64, error)
	createOrUpdateSessionFunc func(userID int64) (string, error)
	getUserIDByTokenFunc      func(token string) (int64, error)
	checkOrderOwnershipFunc   func(orderNumber string, currentUserID int64) (models.OrderCheckResult, error)
	createOrderFunc           func(userId int64, number string, result models.OrderResult) (int, error)
	getOrdersByUserIDFunc     func(userID int64) ([]models.OrderListItem, error)
	getBalanceFunc            func(userID int64) (models.Balance, error)
	accrueBalanceFunc         func(userId int64, amount int64) error
	withdrawFunc              func(orderID string, userID int64, amount int64) error
	getUserWithdrawalsFunc    func(userID int64) ([]models.Withdraw, error)
}

func (m *mockDBRepository) RegisterUser(login, password string) error {
	if m.registerUserFunc != nil {
		return m.registerUserFunc(login, password)
	}
	return nil
}

func (m *mockDBRepository) UserExists(login string) (bool, error) {
	if m.userExistsFunc != nil {
		return m.userExistsFunc(login)
	}
	return false, nil
}

func (m *mockDBRepository) AuthenticateUser(login, password string) (int64, error) {
	if m.authenticateUserFunc != nil {
		return m.authenticateUserFunc(login, password)
	}
	return 0, nil
}

func (m *mockDBRepository) CreateOrUpdateSession(userID int64) (string, error) {
	if m.createOrUpdateSessionFunc != nil {
		return m.createOrUpdateSessionFunc(userID)
	}
	return "test-token", nil
}

func (m *mockDBRepository) GetUserIDByToken(token string) (int64, error) {
	if m.getUserIDByTokenFunc != nil {
		return m.getUserIDByTokenFunc(token)
	}
	return 1, nil
}

func (m *mockDBRepository) CheckOrderOwnership(orderNumber string, currentUserID int64) (models.OrderCheckResult, error) {
	if m.checkOrderOwnershipFunc != nil {
		return m.checkOrderOwnershipFunc(orderNumber, currentUserID)
	}
	return models.OrderCheckResult{}, nil
}

func (m *mockDBRepository) CreateOrder(userId int64, number string, result models.OrderResult) (int, error) {
	if m.createOrderFunc != nil {
		return m.createOrderFunc(userId, number, result)
	}
	return 0, nil
}

func (m *mockDBRepository) GetOrdersByUserID(userID int64) ([]models.OrderListItem, error) {
	if m.getOrdersByUserIDFunc != nil {
		return m.getOrdersByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *mockDBRepository) GetBalance(userID int64) (models.Balance, error) {
	if m.getBalanceFunc != nil {
		return m.getBalanceFunc(userID)
	}
	return models.Balance{}, nil
}

func (m *mockDBRepository) AccrueBalance(userId int64, amount int64) error {
	if m.accrueBalanceFunc != nil {
		return m.accrueBalanceFunc(userId, amount)
	}
	return nil
}

func (m *mockDBRepository) Withdraw(orderID string, userID int64, amount int64) error {
	if m.withdrawFunc != nil {
		return m.withdrawFunc(orderID, userID, amount)
	}
	return nil
}

func (m *mockDBRepository) GetUserWithdrawals(userID int64) ([]models.Withdraw, error) {
	if m.getUserWithdrawalsFunc != nil {
		return m.getUserWithdrawalsFunc(userID)
	}
	return nil, nil
}
