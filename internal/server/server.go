// Package server
package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/handlers"
)

type bodyCaptureResponseWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (w *bodyCaptureResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyCaptureResponseWriter) Write(b []byte) (int, error) {
	// Copy response body to the buffer
	_, _ = w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func formatMaybeJSON(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	var v any
	if err := json.Unmarshal(b, &v); err == nil {
		if pretty, err := json.MarshalIndent(v, "", "  "); err == nil {
			return string(pretty)
		}
	}
	return string(b)
}

func errorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Capture request body
		var reqBody []byte
		if r.Body != nil {
			b, err := io.ReadAll(r.Body)
			if err == nil {
				reqBody = b
			}
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// Capture response
		bw := &bodyCaptureResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(bw, r)

		// Log only on error responses
		if bw.status >= 400 {
			log.Printf("ERROR RESPONSE: %s %s -> %d\nRequestBody: %s\nResponseBody: %s",
				r.Method, r.URL.String(), bw.status, formatMaybeJSON(reqBody), formatMaybeJSON(bw.body.Bytes()))
		}
	})
}

func New(cfg *config.Config) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(errorLogger)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Get("/health", handlers.Health)

	// chatwoot -> codewoot -> codechat
	r.Route("/chatwoot", func(r chi.Router) {
		r.Post("/webhook/{data}", handlers.ChatwootWebhook(cfg))
	})

	// codechat -> codewoot -> chatwoot
	r.Route("/codechat", func(r chi.Router) {
		r.Post("/webhook/{data}", handlers.CodechatWebhook(cfg))
	})

	// TODO: CORS, Auth, Middleware contexto do request (instancia, etc..)
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
