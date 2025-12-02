package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/dto"
	"github.com/sdrvirtual/codewoot/internal/services"
)

// func CodechatWebhook(cfg *config.Config) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var payload dto.CodechatWebhook

// 		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 			http.Error(w, "invalid payload:\n"+err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		relay := services.NewRelayService(cfg)
// 		if err := relay.FromCodechat(payload); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 	}
// }

func CodechatWebhook(cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Content-Type") != "application/json" {
            http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
            return
        }

        // Prevent huge bodies
        r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB

        var payload dto.CodechatWebhook

        // Use Decoder that rejects unknown fields (optional but safer)
        dec := json.NewDecoder(r.Body)

        if err := dec.Decode(&payload); err != nil {
            http.Error(w, "invalid payload:\n"+err.Error(), http.StatusBadRequest)
            return
        }

        relay := services.NewRelayService(cfg)

        if err := relay.FromCodechat(payload); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}
