package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/dto"
	"github.com/sdrvirtual/codewoot/internal/services"
)

func CodechatWebhook(cfg *config.Config, p *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB

		var payload dto.CodechatWebhook

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&payload); err != nil {
			http.Error(w, "invalid payload:\n"+err.Error(), http.StatusBadRequest)
			return
		}

		session := chi.URLParam(r, "session")
		if session == "" {
			http.Error(w, "missing required path param: session", http.StatusBadRequest)
			return
		}

		relay, err := services.NewRelayService(context.WithValue(r.Context(), "session", session), cfg, p, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if err := relay.FromCodechat(payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
