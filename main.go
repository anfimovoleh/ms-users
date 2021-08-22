package app

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/anfimovoleh/ms-users/config"
	"github.com/anfimovoleh/ms-users/server"
)

type App struct {
	config config.Config
	log    *zap.Logger
}

func New(config config.Config) *App {
	return &App{
		config: config,
		log:    config.Log(),
	}
}

func (a *App) Start() error {
	cfg := a.config

	httpCfg := cfg.HTTP()

	router := server.Router(
		cfg,
	)

	serverHost := fmt.Sprintf("%s:%s", httpCfg.Host, httpCfg.Port)
	a.log.With(zap.String("api", "start")).
		Info(fmt.Sprintf("listenig addr =  %s", serverHost))

	httpServer := http.Server{
		Addr:           serverHost,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		return errors.Wrap(err, "failed to start http server")
	}

	return nil
}
