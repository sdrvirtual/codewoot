package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/dto"
	"github.com/sdrvirtual/codewoot/internal/services"
)

// relayTimeout is the maximum time allowed for relaying a message from
// Chatwoot to CodeChat (or vice-versa).  This is intentionally generous
// because audio transcoding + upload can take a while for large files.
const relayTimeout = 180 * time.Second

func ChatwootWebhook(cfg *config.Config, p *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB

		var payload dto.ChatwootWebhook

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

		// Validate the session synchronously so we can return 400 for
		// unknown sessions (this is a fast DB lookup).
		relay, err := services.NewRelayService(r.Context(), cfg, p, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Respond immediately so Chatwoot does not time out waiting for
		// the relay to finish.  The actual work (download, transcode,
		// upload) happens asynchronously in a goroutine with its own
		// independent context and timeout.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), relayTimeout)
			defer cancel()

			relay.SetContext(ctx)

			if err := relay.FromChatwoot(payload); err != nil {
				slog.Error("async relay from chatwoot failed",
					"session", session,
					"error", err.Error(),
				)
			}
		}()
	}
}
