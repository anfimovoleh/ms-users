package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/anfimovoleh/ms-users/db"

	"github.com/go-ozzo/ozzo-validation/v4"

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
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	if err := request.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	token, err := DB(r).GetUserByToken(request.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("Token ID = ", request.Token)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(ErrResponse(http.StatusBadRequest, errors.New("Verification email was already used")))
			return
		}

		log.WithError(err).Error("failed to get user token")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	if token == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ErrResponse(http.StatusBadRequest, errors.New("no such user with provided token address")))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	if err := DB(r).SetUserNewPassword(&db.User{ID: token.UserID, Password: string(hashedPassword)}); err != nil {
		Log(r).WithError(err).Error("failed to update user password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := DB(r).DeleteToken(request.Token); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := DB(r).GetUserByID(token.UserID)
	if err != nil {
		log.WithError(err).Error("failed to get user by id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//notify user about password changing
	if err := EmailClient(r).NewPassword(user.Email); err != nil {
		Log(r).WithField("email_client", "notification").Error("failed to send notification about new password to", token)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
