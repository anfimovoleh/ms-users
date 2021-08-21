package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/anfimovoleh/httperr"

	"github.com/anfimovoleh/ms-users/db"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func (r ResetPasswordRequest) Validate() error {
	return validation.Validate(&r.Email, is.Email, validation.Required)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	log := Log(r).WithField("handler", "reset_password")
	resetPasswordRequest := &ResetPasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(resetPasswordRequest); err != nil {
		httperr.BadRequest(w, errors.New("not valid request body"))
		return
	}

	if err := resetPasswordRequest.Validate(); err != nil {
		httperr.BadRequest(w, err)
		return
	}

	user, err := DB(r).GetUser(resetPasswordRequest.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			httperr.BadRequest(w, errors.New("invalid email address"))
			return
		}

		log.WithError(err).Errorf("failed to get user by email %s", resetPasswordRequest.Email)
		httperr.InternalServerError(w)
		return
	}

	//check if user no exist
	if user == nil {
		//return this by security reasons
		w.WriteHeader(http.StatusAccepted)
		return
	}

	token := uuid.NewString()

	emailToken := &db.Token{
		UserID:     user.ID,
		Token:      token,
		LastSentAt: time.Now(),
	}

	if err := DB(r).CreateToken(emailToken); err != nil {
		Log(r).WithField("db", "email_tokens").
			WithError(err).Error("failed to create token")
		httperr.InternalServerError(w)
		return
	}

	//link to web app new password form
	link := fmt.Sprintf("%s/recovery-password?token=%s", WebApp(r).String(), token)

	//skip err for Email client
	if err := EmailClient(r).Forgot(user.Email, link); err != nil {
		Log(r).Errorf("failed to send forgot password email %s", err)
	}

	w.WriteHeader(http.StatusOK)
}
