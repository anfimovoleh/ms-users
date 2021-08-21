package config

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	"github.com/caarlos0/env"
)

type HTTP struct {
	Host           string `env:"USERS_API_HOST,required"`
	Port           string `env:"USERS_API_PORT,required"`
	SSL            bool   `env:"USERS_API_SSL" envDefault:"false"`
	ServerCertPath string `env:"USERS_API_SERVER_CERT_PATH"`
	ServerKeyPath  string `env:"USERS_API_SERVER_CERT_KEY"`
}

func (h HTTP) URL() (*url.URL, error) {
	if h.SSL {
		resultURL, err := url.Parse(fmt.Sprintf("https://%s:%s", h.Host, h.Port))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse url")
		}

		return resultURL, nil
	}

	resultURL, err := url.Parse(fmt.Sprintf("http://%s:%s", h.Host, h.Port))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	return resultURL, nil
}

func (c *ConfigImpl) HTTP() *HTTP {
	if c.http != nil {
		return c.http
	}

	c.Lock()
	defer c.Unlock()

	http := &HTTP{}
	if err := env.Parse(http); err != nil {
		panic(err)
	}

	c.http = http

	return c.http
}
