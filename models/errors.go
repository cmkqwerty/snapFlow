package models

import "errors"

var (
	ErrEmailTaken = errors.New("models: email address is already taken")
	ErrNotFound   = errors.New("models: resource not found")
)
