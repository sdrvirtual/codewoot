package dto

import (
	"encoding/json"
	"testing"
)

func TestCodechatWebhookUnmarshalLIDFields(t *testing.T) {
	raw := []byte(`{
		"event": "messages.upsert",
		"instance": {
			"name": "galleria",
			"ownerJid": "5521999999999@s.whatsapp.net"
		},
		"data": {
			"id": 123,
			"keyId": "ABC123",
			"keyRemoteJid": "5521984655502@s.whatsapp.net",
			"keyLid": "6846379241653@lid",
			"keyParticipant": "5521984655502@s.whatsapp.net",
			"keyParticipantLid": "6846379241653@lid",
			"keyFromMe": false,
			"pushName": "Cliente",
			"messageType": "conversation",
			"content": {"text": "oi"},
			"messageTimestamp": 1764803328,
			"instanceId": 1,
			"device": "android",
			"isGroup": false
		}
	}`)

	var webhook CodechatWebhook
	if err := json.Unmarshal(raw, &webhook); err != nil {
		t.Fatalf("failed to unmarshal webhook: %v", err)
	}

	if webhook.Data.KeyRemoteJid != "5521984655502@s.whatsapp.net" {
		t.Fatalf("expected keyRemoteJid to be decoded, got %q", webhook.Data.KeyRemoteJid)
	}
	if webhook.Data.KeyLid != "6846379241653@lid" {
		t.Fatalf("expected keyLid to be decoded, got %q", webhook.Data.KeyLid)
	}
	if webhook.Data.KeyParticipant != "5521984655502@s.whatsapp.net" {
		t.Fatalf("expected keyParticipant to be decoded, got %q", webhook.Data.KeyParticipant)
	}
	if webhook.Data.KeyParticipantLid != "6846379241653@lid" {
		t.Fatalf("expected keyParticipantLid to be decoded, got %q", webhook.Data.KeyParticipantLid)
	}

	content, ok := webhook.Data.Content.(CodechatTextContent)
	if !ok {
		t.Fatalf("expected text content, got %T", webhook.Data.Content)
	}
	if content.Text != "oi" {
		t.Fatalf("expected text %q, got %q", "oi", content.Text)
	}
}
