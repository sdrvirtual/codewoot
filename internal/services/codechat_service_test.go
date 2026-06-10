package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sdrvirtual/codewoot/internal/codechat"
	"github.com/sdrvirtual/codewoot/internal/domain"
)

func TestCodechatServiceSendMessageUsesRoutingJID(t *testing.T) {
	var capturedNumber string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/message/sendText/inst" {
			t.Errorf("unexpected request path: %s", r.URL.Path)
		}

		var payload codechat.SendTextParams
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		capturedNumber = payload.Number

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	client, err := codechat.New(ts.URL, "token", codechat.WithInstanceToken("instance-token", "inst"))
	if err != nil {
		t.Fatalf("failed to create codechat client: %v", err)
	}

	service := &CodechatService{client: client}
	contact := domain.ContactInfo{
		Phone:      "5521984655502",
		RoutingJID: "6846379241653@lid",
	}
	message := CodechatClientMessage{Text: "oi"}

	if err := service.SendMessage(t.Context(), contact, message); err != nil {
		t.Fatalf("SendMessage returned error: %v", err)
	}

	if capturedNumber != contact.RoutingJID {
		t.Fatalf("expected number %q, got %q", contact.RoutingJID, capturedNumber)
	}
}

func TestCodechatServiceSendMessageFallsBackToPhone(t *testing.T) {
	var capturedNumber string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload codechat.SendTextParams
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		capturedNumber = payload.Number

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	client, err := codechat.New(ts.URL, "token", codechat.WithInstanceToken("instance-token", "inst"))
	if err != nil {
		t.Fatalf("failed to create codechat client: %v", err)
	}

	service := &CodechatService{client: client}
	contact := domain.ContactInfo{Phone: "5521984655502"}
	message := CodechatClientMessage{Text: "oi"}

	if err := service.SendMessage(t.Context(), contact, message); err != nil {
		t.Fatalf("SendMessage returned error: %v", err)
	}

	if capturedNumber != contact.Phone {
		t.Fatalf("expected number %q, got %q", contact.Phone, capturedNumber)
	}
}

func TestCodechatServiceResolveRoutingNumberUsesLID(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/whatsappNumbers/inst" {
			t.Errorf("unexpected request path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"jid":"5521984655502@s.whatsapp.net","exists":true,"lid":"6846379241653@lid"}]`))
	}))
	defer ts.Close()

	client, err := codechat.New(ts.URL, "token", codechat.WithInstanceToken("instance-token", "inst"))
	if err != nil {
		t.Fatalf("failed to create codechat client: %v", err)
	}

	service := &CodechatService{client: client}
	route, err := service.ResolveRoutingNumber(t.Context(), "5521984655502")
	if err != nil {
		t.Fatalf("ResolveRoutingNumber returned error: %v", err)
	}

	if route.Number != "6846379241653@lid" {
		t.Fatalf("expected routing number to use lid, got %q", route.Number)
	}
	if route.JID != "5521984655502@s.whatsapp.net" || route.LID != "6846379241653@lid" || !route.Exists {
		t.Fatalf("unexpected route: %#v", route)
	}
}

func TestCodechatServiceResolveRoutingNumberFallsBackToPhone(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"jid":"5521984655502@s.whatsapp.net","exists":true}]`))
	}))
	defer ts.Close()

	client, err := codechat.New(ts.URL, "token", codechat.WithInstanceToken("instance-token", "inst"))
	if err != nil {
		t.Fatalf("failed to create codechat client: %v", err)
	}

	service := &CodechatService{client: client}
	route, err := service.ResolveRoutingNumber(t.Context(), "5521984655502")
	if err != nil {
		t.Fatalf("ResolveRoutingNumber returned error: %v", err)
	}

	if route.Number != "5521984655502" {
		t.Fatalf("expected routing number to fall back to phone, got %q", route.Number)
	}
}
