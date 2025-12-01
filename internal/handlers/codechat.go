package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/dto"
	"github.com/sdrvirtual/codewoot/internal/services"
)

func CodechatWebhook(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload dto.CodechatWebhook

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		relay := services.NewRelayService(cfg)
		relay.FromCodechat(payload)

	}
}
