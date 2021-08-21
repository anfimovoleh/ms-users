package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-ozzo/ozzo-validation/v4"

	"github.com/dgrijalva/jwt-go"

	"github.com/json-iterator/go"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c LoginRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required),
		validation.Field(&c.Password, validation.Required),
	)
}

type LoginResponse struct {
	Token string `json:"token"`
}

const tokenExpirationDuration = time.Hour

func Login(w http.ResponseWriter, r *http.Request) {
	loginRequest := &LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	if err := loginRequest.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	user, err := DB(r).GetUser(loginRequest.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(ErrResponse(http.StatusUnauthorized, ErrInvalidEmailOrPassword))
			return
		}

		Log(r).WithError(err).Error("failed to get user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(ErrResponse(http.StatusUnauthorized, ErrInvalidEmailOrPassword))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(ErrResponse(http.StatusUnauthorized, ErrInvalidEmailOrPassword))
		return
	}
	_, token, err := JWT(r).Encode(
		jwt.MapClaims{
			"id":  user.ID,
			"exp": time.Now().Add(tokenExpirationDuration).Unix(),
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(ErrResponse(http.StatusBadRequest, err))
		return
	}

	result := LoginResponse{
		Token: token,
	}

	serializer := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 "json",
	}.Froze()
	response, err := serializer.Marshal(result)
	if err != nil {
		Log(r).WithField("response", "LoginResponseSerialize").Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write(response)
}
