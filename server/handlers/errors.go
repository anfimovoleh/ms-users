package handlers

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyRequestToken      = errors.New("request token should not be empty")
	ErrEmptyPassword          = errors.New("password should not be empty")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

func ErrResponse(code int, err error) []byte {
	return []byte(fmt.Sprintf(`{"code": %d, "error": "%s"}`, code, err.Error()))
}
