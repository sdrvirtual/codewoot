package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/dto"
	"github.com/sdrvirtual/codewoot/internal/services"
)

func ChatwootWebhook(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload dto.ChatwootWebhook

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		relay := services.NewRelayService(cfg)
		relay.FromChatwoot(payload)

		// s, _ := json.MarshalIndent(payload, "", "\t")
		// fmt.Println(string(s))

		w.WriteHeader(http.StatusOK)
	}
}
