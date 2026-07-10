package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	errorsInternal "github.com/postman17/gofermart/internal/errors"
	hash "github.com/postman17/gofermart/internal/hash"
)

func (d *DBStorage) RegisterUser(ctx context.Context, login string, password string) error {
	exists, err := d.UserExists(ctx, login)
	if err != nil {
		return fmt.Errorf("user exists error '%s'", login)
	}

	if exists {
		return fmt.Errorf("%w: %s", errorsInternal.ErrUserAlreadyExists, login)
	}

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := "INSERT INTO users (login, password) VALUES ($1, $2)"
	_, err = d.db.ExecContext(ctx, query, login, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (d *DBStorage) UserExists(ctx context.Context, login string) (bool, error) {
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE login = $1)"
	err := d.db.QueryRowContext(ctx, query, login).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (d *DBStorage) AuthenticateUser(ctx context.Context, login string, password string) (int64, error) {
	var (
		userID     int64
		storedHash string
	)

	query := "SELECT id, password FROM users WHERE login = $1"
	err := d.db.QueryRowContext(ctx, query, login).Scan(&userID, &storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to query user: %w", err)
	}

	if !hash.CheckPassword(password, storedHash) {
		return 0, nil
	}

	return userID, nil
}

func (d *DBStorage) CreateOrUpdateSession(ctx context.Context, userID int64) (string, error) {
	token, err := generateRandomToken(20)
	if err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	query := `
		INSERT INTO sessions (user_id, token, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			token = EXCLUDED.token,
			updated_at = NOW()
		RETURNING token;
	`

	var returnedToken string
	err = d.db.QueryRowContext(ctx, query, userID, token).Scan(&returnedToken)
	if err != nil {
		return "", fmt.Errorf("failed to upsert session in database: %w", err)
	}

	return returnedToken, nil
}

func generateRandomToken(length int) (string, error) {
	// 1 байт в hex-кодировке превращается в 2 символа.
	// Поэтому делим длину пополам. Для 20 символов нам нужно 10 байт.
	bytesCount := length / 2
	bytes := make([]byte, bytesCount)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func (d *DBStorage) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	var userID int64

	query := `SELECT user_id FROM sessions WHERE token = $1 LIMIT 1;`
	err := d.db.QueryRowContext(ctx, query, token).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errorsInternal.ErrSessionNotFound
		}
		return 0, err
	}

	return userID, nil
}
