package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/caarlos0/env"
)

type HTTP struct {
	Host            string        `env:"USERS_API_HOST,required"`
	Port            string        `env:"USERS_API_PORT,required"`
	ReqDurThreshold time.Duration `env:"USERS_HTTP_REQ_DUR_THRESHOLD" envDefault:"5s"`
}

func (h HTTP) URL() (*url.URL, error) {
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

	http := &HTTP{}
	if err := env.Parse(http); err != nil {
		panic(err)
	}

	c.http = http

	return c.http
}
