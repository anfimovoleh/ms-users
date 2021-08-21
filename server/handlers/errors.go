package handlers

import (
	"errors"
)

var (
	ErrEmptyRequestToken      = errors.New("request token should not be empty")
	ErrEmptyPassword          = errors.New("password should not be empty")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)