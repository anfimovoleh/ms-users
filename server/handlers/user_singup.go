package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/anfimovoleh/httperr"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/anfimovoleh/ms-users/db"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	DateOfBirth string `json:"date_of_birth"`
}

func (u SignupRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Name, validation.Required),
		validation.Field(&u.Phone, validation.Required),
	)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	log := Log(r).WithField("handler", "user_signup")
	signupRequest := &SignupRequest{}
	if err := json.NewDecoder(r.Body).Decode(signupRequest); err != nil {
		httperr.BadRequest(w, err)
		return
	}

	if err := signupRequest.Validate(); err != nil {
		httperr.BadRequest(w, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupRequest.Password), 8)
	if err != nil {
		httperr.BadRequest(w, err)
		return
	}

	createdUser, err := DB(r).GetUser(signupRequest.Email)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			createdUser = nil
		default:
			log.WithError(err).Errorf("failed to get user with email = %s", signupRequest.Email)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if createdUser != nil {
		httperr.BadRequest(w, errors.New("user with this email is already registered"))
		return
	}

	dbUser := &db.User{
		Name:        signupRequest.Name,
		Email:       signupRequest.Email,
		Password:    string(hashedPassword),
		Phone:       signupRequest.Phone,
		DateOfBirth: signupRequest.DateOfBirth,
	}

	if err := DB(r).CreateUser(dbUser); err != nil {
		Log(r).WithField("table", "users").Error(err)
		httperr.InternalServerError(w)
		return
	}

	user, err := DB(r).GetUser(dbUser.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			httperr.BadRequest(w, errors.New("invalid email address"))
			return
		}

		log.WithError(err).Errorf("failed to get user by email %s", dbUser.Email)
		httperr.InternalServerError(w)
		return
	}

	token := uuid.NewString()
	confirmToken := &db.Token{
		UserID:     user.ID,
		Token:      token,
		LastSentAt: time.Now(),
	}

	if err := DB(r).CreateToken(confirmToken); err != nil {
		Log(r).WithError(err).Error("failed to create token")
		httperr.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
