package models

import (
	"database/sql"
	"fmt"
)

type Session struct {
	ID        int
	UserID    int
	Token     string // Only set while creating new session. Otherwise, it is blank.
	TokenHash string
}

type SessionService struct {
	DB           *sql.DB
	TokenManager TokenManager
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	token, tokenHash, err := ss.TokenManager.New()
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: tokenHash,
	}

	row := ss.DB.QueryRow(`
	INSERT INTO sessions (user_id, token_hash)
	VALUES ($1, $2) ON CONFLICT (user_id) DO
	UPDATE
	SET token_hash = $2
	RETURNING id;`, session.UserID, session.TokenHash)
	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.TokenManager.hash(token)

	var user User
	row := ss.DB.QueryRow(`
	SELECT users.id,
	users.email,
	users.password_hash
	FROM sessions
	JOIN users ON users.id = sessions.user_id
	WHERE sessions.token_hash = $1;`, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}

	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.TokenManager.hash(token)

	_, err := ss.DB.Exec(`
	DELETE FROM sessions
	WHERE token_hash = $1;`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
