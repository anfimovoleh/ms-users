package utils

import (
	"context"
	"fmt"

	"github.com/anfimovoleh/ms-users/db"

	"github.com/pkg/errors"

	"github.com/go-chi/jwtauth/v5"
)

func ErrResponse(code int, err error) []byte {
	return []byte(fmt.Sprintf(`{"code": %d, "error": "%s"}`, code, err.Error()))
}

var (
	ErrInvalidSession         = errors.New("invalid session")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

func User(ctx context.Context, db *db.DB) (*db.User, string, error) {
	token, claims, _ := jwtauth.FromContext(ctx)
	userEmail, ok := claims["user_email"].(string)
	if !ok {
		return nil, "", ErrInvalidSession
	}

	dbUser, err := db.GetUser(userEmail)
	if err != nil {
		return nil, "", err
	}

	if dbUser == nil {
		return nil, "", ErrInvalidEmailOrPassword
	}

	return dbUser, token.Subject(), nil
}
