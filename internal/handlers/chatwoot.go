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
			http.Error(w, "invalid payload:\n"+err.Error(), http.StatusBadRequest)
			return
		}

		// a, _ := io.ReadAll(r.Body)
		// s, _ := json.MarshalIndent(json.RawMessage(a), "", "  ")
		// fmt.Println("------- FromChatwoot ----------\n", string(s))

		relay := services.NewRelayService(cfg)
		if err := relay.FromChatwoot(payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
