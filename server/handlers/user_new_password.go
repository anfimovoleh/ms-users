package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/anfimovoleh/httperr"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/anfimovoleh/ms-users/db"

	"golang.org/x/crypto/bcrypt"
)

type NewPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (n NewPasswordRequest) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(&n.Password, validation.Required),
		validation.Field(&n.Token, validation.Required),
	)
}

type NewPasswordHandler struct {
	log *zap.Logger
}

func NewNewPasswordHandler(log *zap.Logger) *NewPasswordHandler {
	return &NewPasswordHandler{log: log}
}

func (h NewPasswordHandler) Handle(w http.ResponseWriter, r *http.Request) {
	request := &NewPasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		httperr.BadRequest(w, err)
		return
	}

	if err := request.Validate(); err != nil {
		httperr.BadRequest(w, err)
		return
	}

	token, err := DB(r).GetUserByToken(request.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			httperr.BadRequest(w, errors.New("Verification email was already used"))
			return
		}

		h.log.With(
			zap.Error(err),
		).Error("failed to get user token")
		httperr.InternalServerError(w)
		return
	}

	if token == nil {
		httperr.BadRequest(w, errors.New("no such user with provided token address"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 8)
	if err != nil {
		httperr.BadRequest(w, err)
		return
	}

	if err := DB(r).SetUserNewPassword(&db.User{ID: token.UserID, Password: string(hashedPassword)}); err != nil {
		h.log.With(
			zap.Error(err),
		).Error("failed to update user password")
		httperr.InternalServerError(w)
		return
	}

	if err := DB(r).DeleteToken(request.Token); err != nil {
		httperr.InternalServerError(w)
		return
	}

	user, err := DB(r).GetUserByID(token.UserID)
	if err != nil {
		h.log.With(
			zap.Error(err),
		).Error("failed to get user by id")
		httperr.InternalServerError(w)
		return
	}

	//notify user about password changing
	if err := EmailClient(r).NewPassword(user.Email); err != nil {
		h.log.With(
			zap.String("email_client", "notification"),
			zap.Any("token", token),
		).Error("failed to send notification about new password")
		httperr.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
