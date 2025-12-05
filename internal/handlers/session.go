package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdrvirtual/codewoot/internal/config"
)

func CreateSession(cfg *config.Config, p *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Api-Key") != cfg.Authorization.Key {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
	}
}
