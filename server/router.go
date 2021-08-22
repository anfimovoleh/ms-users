package server

import (
	"net/http"

	chiwares "github.com/anfimovoleh/go-chi-middlewares"

	"github.com/anfimovoleh/ms-users/config"

	"github.com/anfimovoleh/ms-users/server/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func Router(
	cfg config.Config,
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

	url, err := cfg.HTTP().URL()
	if err != nil {
		cfg.Log().Fatal("failed to get URL")
	}

	router.Use(
		cors.Handler,
		chiwares.Logger(cfg.Log(), cfg.HTTP().ReqDurThreshold),
		chiwares.Ctx(
			handlers.CtxHTTP(url),
			handlers.CtxEmailClient(cfg.EmailClient()),
			handlers.CtxWebApp(cfg.WebsiteURL()),
			handlers.CtxDB(cfg.DB()),
			handlers.CtxJWT(cfg.JWT()),
		),
	)

	router.Route("/user", func(router chi.Router) {
		router.Post("/login", handlers.NewLoginHandler(cfg.Log()).Handle)
		router.Post("/signup", handlers.NewSignupHandler(cfg.Log()).Handle)
		router.Put("/new_password", handlers.NewNewPasswordHandler(cfg.Log()).Handle)
		router.Post("/reset_password", handlers.NewResetPasswordHandler(cfg.Log()).Handle)
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"api":"ms-users"}`))
	})

	return router
}
