package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	models "github.com/postman17/gofermart/internal/model"
)

func (d *DBStorage) CheckOrderOwnership(ctx context.Context, orderNumber string, currentUserID int64) (models.OrderCheckResult, error) {
	var dbUserID int64

	query := `SELECT user_id FROM orders WHERE number = $1 LIMIT 1;`

	err := d.db.QueryRowContext(ctx, query, orderNumber).Scan(&dbUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.OrderCheckResult{Exists: false, IsSameUser: false}, nil
		}
		return models.OrderCheckResult{}, fmt.Errorf("failed to check order ownership: %w", err)
	}

	isSame := (dbUserID == currentUserID)
	return models.OrderCheckResult{Exists: true, IsSameUser: isSame}, nil
}

func (d *DBStorage) CreateOrder(ctx context.Context, userId int64, number string, result models.OrderResult) (int, error) {
	var lastInsertID int
	now := time.Now().UTC()

	query := `
		INSERT INTO orders (user_id, number, status, accrual, uploaded_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	err := d.db.QueryRowContext(ctx, query, userId, number, result.Status, result.Accrual, now).Scan(&lastInsertID)

	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

func (d *DBStorage) GetOrdersByUserID(ctx context.Context, userID int64) ([]models.OrderListItem, error) {
	query := `
		SELECT id, user_id, number, status, accrual, uploaded_at, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrderListItem

	for rows.Next() {
		var o models.OrderListItem
		err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.Number,
			&o.Status,
			&o.Accrual,
			&o.UploadedAt,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
