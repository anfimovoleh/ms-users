package config

import (
	"github.com/caarlos0/env"
	"net/url"

	"github.com/pkg/errors"
)

type Website struct {
	URL string `env:"USERS_WEBSITE_URL,required"`
}

func (c *ConfigImpl) WebsiteURL() *url.URL {
	if c.webApp != nil {
		return c.webApp
	}

	c.Lock()
	defer c.Unlock()

	webApp := &Website{}
	if err := env.Parse(webApp); err != nil {
		panic(err)
	}

	url, err := url.Parse(webApp.URL)
	if err != nil {
		panic(errors.Wrap(err, "invalid url"))
	}

	c.webApp = url

	return c.webApp
}
