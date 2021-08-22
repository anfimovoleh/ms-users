package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/anfimovoleh/httperr"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/dgrijalva/jwt-go"

	jsoniter "github.com/json-iterator/go"
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

type LoginHandler struct {
	log *zap.Logger
}

func NewLoginHandler(log *zap.Logger) *LoginHandler {
	return &LoginHandler{log: log}
}

func (h LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	loginRequest := &LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(loginRequest); err != nil {
		httperr.BadRequest(w, err)
		return
	}

	if err := loginRequest.Validate(); err != nil {
		httperr.BadRequest(w, err)
		return
	}

	user, err := DB(r).GetUser(loginRequest.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			httperr.BadRequest(w, err)
			return
		}

		h.log.With(
			zap.String("email", loginRequest.Email),
			zap.Error(err),
		).Error("failed to get user")
		httperr.InternalServerError(w)
		return
	}

	if user == nil {
		httperr.ErrResponse(w, http.StatusUnauthorized, ErrInvalidEmailOrPassword)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		httperr.ErrResponse(w, http.StatusUnauthorized, ErrInvalidEmailOrPassword)
		return
	}
	_, token, err := JWT(r).Encode(
		jwt.MapClaims{
			"id":  user.ID,
			"exp": time.Now().Add(tokenExpirationDuration).Unix(),
		},
	)
	if err != nil {
		httperr.BadRequest(w, err)
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
		h.log.With(
			zap.Error(err),
		).Error("failed to serialize response")
		httperr.InternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write(response)
}
