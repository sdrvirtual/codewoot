package codechat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhatsappNumbers(t *testing.T) {
	var capturedNumbers []string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/whatsappNumbers/inst" {
			t.Fatalf("expected whatsappNumbers path, got %q", r.URL.Path)
		}

		var payload WhatsappNumbersParams
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		capturedNumbers = payload.Numbers

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"jid":"5521984655502@s.whatsapp.net","exists":true,"lid":"6846379241653@lid"}]`))
	}))
	defer ts.Close()

	client, err := New(ts.URL, "token", WithInstanceToken("instance-token", "inst"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	resp, err := client.WhatsappNumbers(t.Context(), []string{"5521984655502"})
	if err != nil {
		t.Fatalf("WhatsappNumbers returned error: %v", err)
	}

	if len(capturedNumbers) != 1 || capturedNumbers[0] != "5521984655502" {
		t.Fatalf("unexpected request numbers: %#v", capturedNumbers)
	}
	if len(resp) != 1 {
		t.Fatalf("expected one response, got %d", len(resp))
	}
	if resp[0].JID != "5521984655502@s.whatsapp.net" || !resp[0].Exists || resp[0].LID != "6846379241653@lid" {
		t.Fatalf("unexpected response: %#v", resp[0])
	}
}
