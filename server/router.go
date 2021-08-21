package server

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/jwtauth"

	"github.com/anfimovoleh/ms-users/db"
	"github.com/anfimovoleh/ms-users/email"
	"github.com/anfimovoleh/ms-users/server/handlers"
	"github.com/anfimovoleh/ms-users/server/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

const durationThreshold = time.Second * 10

func Router(
	log *logrus.Entry,
	emailClient email.Client,
	url *url.URL,
	webApp *url.URL,
	db *db.DB,
	jwtAuth *jwtauth.JWTAuth,
) chi.Router {
	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*", "https://localhost:3000"},
		AllowedMethods:   []string{"*", "GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-auth"},
		ExposedHeaders:   []string{"*", "Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	router.Use(
		cors.Handler,
		middlewares.Logger(log, durationThreshold),
		middlewares.Ctx(
			handlers.CtxLog(log),
			handlers.CtxHTTP(url),
			handlers.CtxEmailClient(emailClient),
			handlers.CtxWebApp(webApp),
			handlers.CtxDB(db),
			handlers.CtxJWT(jwtAuth),
		),
	)

	router.Route("/user", func(router chi.Router) {
		router.Post("/login", handlers.Login)
		router.Post("/signup", handlers.Signup)
		router.Put("/new_password", handlers.NewPassword)
		router.Post("/reset_password", handlers.ResetPassword)
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"api":"ms-users"}`))
	})

	return router
}
