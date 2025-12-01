package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/handlers"
)

func New(cfg *config.Config) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

    r.Get("/health", handlers.Health)
	// chatwoot -> codewoot -> codechat
	r.Route("/chatwoot", func (r chi.Router) {
		r.Post("/webhook/{data}", handlers.ChatwootWebhook(cfg))
	})
	// codechat -> codewoot -> chatwoot
	r.Route("/codechat", func (r chi.Router) {
		r.Post("/webhook/{data}", handlers.CodechatWebhook(cfg))
	})

    // TODO: CORS, Auth, Middleware contexto do request (instancia, etc..)
    addr := cfg.Server.Host + ":" + cfg.Server.Port
    return &http.Server{
        Addr:    addr,
        Handler: r,
    }
}
