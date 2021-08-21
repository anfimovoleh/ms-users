package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

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

func NewPassword(w http.ResponseWriter, r *http.Request) {
	log := Log(r).WithField("db", "email_tokens")

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
			log.Debug("Token ID = ", request.Token)
			httperr.BadRequest(w, errors.New("Verification email was already used"))
			return
		}

		log.WithError(err).Error("failed to get user token")
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
		Log(r).WithError(err).Error("failed to update user password")
		httperr.InternalServerError(w)
		return
	}

	if err := DB(r).DeleteToken(request.Token); err != nil {
		httperr.InternalServerError(w)
		return
	}

	user, err := DB(r).GetUserByID(token.UserID)
	if err != nil {
		log.WithError(err).Error("failed to get user by id")
		httperr.InternalServerError(w)
		return
	}

	//notify user about password changing
	if err := EmailClient(r).NewPassword(user.Email); err != nil {
		Log(r).WithField("email_client", "notification").Error("failed to send notification about new password to", token)
		httperr.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
