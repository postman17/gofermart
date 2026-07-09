package repository

import (
	"context"
	"fmt"

	errorInternal "github.com/postman17/gofermart/internal/errors"
	models "github.com/postman17/gofermart/internal/model"
)

func (d *DBStorage) Withdraw(ctx context.Context, orderID string, userID int64, amount int64) error {
	if amount <= 0 {
		return errorInternal.ErrInvalidAmount
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE user_balances 
		SET 
			current = current - $1::NUMERIC, 
			withdrawn = withdrawn + $1::NUMERIC,
			updated_at = NOW()
		WHERE user_id = $2 AND current >= $1::NUMERIC`

	result, err := tx.ExecContext(ctx, query, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to execute debit query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errorInternal.ErrInsufficientFunds
	}

	insertWithdrawQuery := `
		INSERT INTO withdraw (user_id, "order", total_sum) 
		VALUES ($1, $2, $3)`

	_, err = tx.ExecContext(ctx, insertWithdrawQuery, userID, orderID, amount)
	if err != nil {
		return fmt.Errorf("failed to insert withdraw record: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

func (d *DBStorage) GetUserWithdrawals(ctx context.Context, userID int64) ([]models.Withdraw, error) {
	query := `
		SELECT id, user_id, "order", total_sum, created_at 
		FROM withdraw 
		WHERE user_id = $1 
		ORDER BY created_at DESC`

	rows, err := d.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query withdrawals: %w", err)
	}
	defer rows.Close()

	withdrawals := make([]models.Withdraw, 0)

	for rows.Next() {
		var w models.WithdrawalsModel
		err := rows.Scan(&w.ID, &w.UserID, &w.Order, &w.TotalSum, &w.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan withdraw row: %w", err)
		}
		withdrawals = append(withdrawals, models.Withdraw{Order: w.Order, TotalSum: w.TotalSum})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return withdrawals, nil
}
