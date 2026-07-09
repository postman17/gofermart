package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	models "github.com/postman17/gofermart/internal/model"
)

func (d *DBStorage) GetBalance(ctx context.Context, userID int64) (models.Balance, error) {
	query := `
		SELECT current, withdrawn 
		FROM user_balances 
		WHERE user_id = $1;
	`

	var b models.Balance

	err := d.db.QueryRowContext(ctx, query, userID).Scan(&b.Current, &b.Withdrawn)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Balance{
				Current:   0.00,
				Withdrawn: 0.00,
			}, nil
		}
		return models.Balance{}, err
	}

	return b, nil
}

func (d *DBStorage) AccrueBalance(ctx context.Context, userId int64, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than 0: %d", amount)
	}

	query := `
		INSERT INTO user_balances (user_id, current, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			current = user_balances.current + EXCLUDED.current,
			updated_at = NOW();
	`

	_, err := d.db.ExecContext(ctx, query, userId, amount)
	if err != nil {
		return fmt.Errorf("error in accrual %d: %w", userId, err)
	}

	return nil
}
