package app

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/anfimovoleh/ms-users/config"
	"github.com/anfimovoleh/ms-users/server"

	"github.com/sirupsen/logrus"
)

type App struct {
	config config.Config
	log    *logrus.Entry
}

func New(config config.Config) *App {
	return &App{
		config: config,
		log:    config.Log(),
	}
}

func (a *App) Start() error {
	conf := a.config

	httpConfiguration := conf.HTTP()

	url, err := httpConfiguration.URL()
	if err != nil {
		return err
	}

	router := server.Router(
		conf.Log(),
		conf.EmailClient(),
		url,
		conf.WebsiteURL(),
		conf.DB(),
		conf.JWT(),
	)

	serverHost := fmt.Sprintf("%s:%s", httpConfiguration.Host, httpConfiguration.Port)
	a.log.WithField("api", "start").
		Info(fmt.Sprintf("listenig addr =  %s, tls = %v", serverHost, httpConfiguration.SSL))

	httpServer := http.Server{
		Addr:           serverHost,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	switch httpConfiguration.SSL {
	case true:
		tlsConfig := &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

				// Best disabled, as they don't provide Forward Secrecy,
				// but might be necessary for some clients
				// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519, // Go 1.8 only
			},
			InsecureSkipVerify: true,
		}

		httpServer.TLSConfig = tlsConfig
		if err := httpServer.ListenAndServeTLS(httpConfiguration.ServerCertPath, httpConfiguration.ServerKeyPath); err != nil {
			return errors.Wrap(err, "failed to start https server")
		}

	default:
		if err := httpServer.ListenAndServe(); err != nil {
			return errors.Wrap(err, "failed to start http server")
		}
	}

	return nil
}
