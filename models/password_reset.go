package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID        int
	UserID    int
	Token     string // Only set while creating new PasswordReset. Otherwise, it is blank.
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB           *sql.DB
	TokenManager TokenManager
	Duration     time.Duration // Amount of time that PasswordReset is valid for.
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)

	var userID int

	row := service.DB.QueryRow(`
	SELECT id FROM users WHERE email = $1;`, email)

	err := row.Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	token, tokenHash, err := service.TokenManager.New()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	duration := service.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	passwordReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(duration),
	}

	row = service.DB.QueryRow(`
	INSERT INTO password_reset (user_id, token_hash, expires_at)
	VALUES ($1, $2, $3) ON CONFLICT (user_id) DO
	UPDATE
	SET token_hash = $2, expires_at = $3
	RETURNING id;`, passwordReset.UserID, passwordReset.TokenHash, passwordReset.ExpiresAt)

	err = row.Scan(&passwordReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &passwordReset, nil
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordResetService.Consume")
}
