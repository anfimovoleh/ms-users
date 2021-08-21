package config

import (
	"github.com/caarlos0/env"
	"github.com/go-chi/jwtauth"
)

type Authentication struct {
	VerifyKey string `env:"USERS_AUTHENTICATION_SECRET,required"`
	Algorithm string `env:"USERS_AUTHENTICATION_ALGORITHM" envDefault:"HS256"`
}

func (jwt *Authentication) GetJWTEntry() *jwtauth.JWTAuth {
	return jwtauth.New(jwt.Algorithm, []byte(jwt.VerifyKey), nil)
}

func (c *ConfigImpl) JWT() *jwtauth.JWTAuth {
	if c.jwt != nil {
		return c.jwt
	}

	c.Lock()
	defer c.Unlock()

	jwt := &Authentication{}
	if err := env.Parse(jwt); err != nil {
		panic(err)
	}

	c.jwt = jwt.GetJWTEntry()

	return c.jwt
}
