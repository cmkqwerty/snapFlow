package models

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/cmkqwerty/snapFlow/rand"
)

const (
	MinBytesPerToken = 32
)

type TokenManager struct {
	BytesPerToken int
}

func (tm TokenManager) New() (token, tokenHash string, err error) {
	bytesPerToken := tm.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err = rand.String(bytesPerToken)
	if err != nil {
		return "", "", fmt.Errorf("token create: %w", err)
	}

	return token, tm.hash(token), nil
}

func (tm TokenManager) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))

	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
