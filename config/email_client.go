package config

import (
	"github.com/anfimovoleh/ms-users/email"
	"github.com/caarlos0/env"
)

type EmailClient struct {
	EmailAddress string `env:"USERS_EMAIL_ADDRESS,required"`
	Password     string `env:"USERS_EMAIL_PASSWORD,required"`
	Host         string `env:"USERS_SMTP_SERVER_HOST" envDefault:"smtp.gmail.com"`
	Port         int    `env:"USERS_SMTP_SERVER_PORT" envDefault:"465"`
}

func (c *ConfigImpl) EmailClient() *email.ClientImpl {
	if c.email != nil {
		return c.email
	}

	c.Lock()
	defer c.Unlock()

	emailClient := &EmailClient{}
	if err := env.Parse(emailClient); err != nil {
		panic(err)
	}

	c.email = email.New(
		emailClient.EmailAddress,
		emailClient.Password,
		emailClient.Host,
		emailClient.Port,
	)
	return c.email
}
