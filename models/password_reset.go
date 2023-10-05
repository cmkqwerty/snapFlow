package models

import (
	"database/sql"
	"fmt"
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

func (service *PasswordResetService) Create(mail string) (*PasswordReset, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordResetService.Create")
}

func (service *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("TODO: Implement PasswordResetService.Consume")
}
