package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sdrvirtual/codewoot/internal/codechat"
	"github.com/sdrvirtual/codewoot/internal/dto"
)

func TestRelayServiceFromChatwoot_NormalizesProblematicMarkdownLinksForOutgoingMessages(t *testing.T) {
	const mapsURL = "https://maps.app.goo.gl/UYZvebsBMEzk7cqk8"

	tests := []struct {
		name       string
		senderType string
		input      string
		want       string
	}{
		{
			name:       "agent bold autolink in sentence",
			senderType: "User",
			input:      "Segue o endereco: **<" + mapsURL + ">** ",
			want:       "Segue o endereco: " + mapsURL,
		},
		{
			name:       "robot bold autolink in sentence",
			senderType: "AgentBot",
			input:      "Segue o endereco: **<" + mapsURL + ">** ",
			want:       "Segue o endereco: " + mapsURL,
		},
		{
			name:       "robot bare bold link",
			senderType: "AgentBot",
			input:      "**" + mapsURL + "**",
			want:       mapsURL,
		},
		{
			name:       "robot dangling bold marker after link",
			senderType: "AgentBot",
			input:      mapsURL + "**",
			want:       mapsURL,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			capturedText := sendChatwootTextThroughRelay(t, tc.input, tc.senderType)
			if capturedText != tc.want {
				t.Fatalf("sent text mismatch\nwant: %q\n got: %q", tc.want, capturedText)
			}
		})
	}
}

func TestRelayServiceFromChatwoot_LeavesSafeOutgoingMessagesUnchanged(t *testing.T) {
	const mapsURL = "https://maps.app.goo.gl/UYZvebsBMEzk7cqk8"

	tests := []struct {
		name       string
		senderType string
		input      string
	}{
		{name: "agent plain text", senderType: "User", input: "Ola, tudo bem?"},
		{name: "robot plain text", senderType: "AgentBot", input: "Ola, tudo bem?"},
		{name: "robot plain URL", senderType: "AgentBot", input: mapsURL},
		{name: "robot non-link bold", senderType: "AgentBot", input: "**Atencao**"},
		{name: "robot chatwoot italic text", senderType: "AgentBot", input: "*Teste*"},
		{name: "robot whatsapp native bold URL", senderType: "AgentBot", input: "*" + mapsURL + "*"},
		{name: "robot bold markers separated from URL", senderType: "AgentBot", input: "Teste link: ** " + mapsURL + " **"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			capturedText := sendChatwootTextThroughRelay(t, tc.input, tc.senderType)
			if capturedText != tc.input {
				t.Fatalf("message should not be normalized\nwant: %q\n got: %q", tc.input, capturedText)
			}
		})
	}
}

func TestRelayServiceFromChatwoot_ContractTextWebhookToCodechat(t *testing.T) {
	const mapsURL = "https://maps.app.goo.gl/UYZvebsBMEzk7cqk8"

	codechatSvc, requests := newCodechatServiceRecorder(t)
	relay := newRelayServiceForCodechatTest(codechatSvc)
	payload := mustDecodeChatwootWebhook(t, chatwootWebhookJSON("Segue endereco: **<"+mapsURL+">** ", ""))

	if err := relay.FromChatwoot(payload); err != nil {
		t.Fatalf("FromChatwoot returned error: %v", err)
	}

	req := requireSingleCodechatRequest(t, requests)
	if req.Method != http.MethodPost {
		t.Fatalf("expected POST, got %s", req.Method)
	}
	if req.Path != "/message/sendText/test-instance" {
		t.Fatalf("expected sendText URL, got %q", req.Path)
	}
	if got := req.Header.Get("apikey"); got != "global-token" {
		t.Fatalf("expected global apikey, got %q", got)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer instance-token" {
		t.Fatalf("expected instance authorization, got %q", got)
	}

	var body struct {
		Number      string `json:"number"`
		TextMessage struct {
			Text string `json:"text"`
		} `json:"textMessage"`
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		t.Fatalf("failed to decode CodeChat request body: %v", err)
	}
	if body.Number != "5549991158078" {
		t.Fatalf("expected normalized phone number, got %q", body.Number)
	}
	wantText := "Segue endereco: " + mapsURL
	if body.TextMessage.Text != wantText {
		t.Fatalf("expected normalized text %q, got %q", wantText, body.TextMessage.Text)
	}
}

func TestRelayServiceFromChatwoot_ContractImageCaptionToCodechat(t *testing.T) {
	const mapsURL = "https://maps.app.goo.gl/UYZvebsBMEzk7cqk8"

	var mediaRequests []string
	mediaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mediaRequests = append(mediaRequests, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(mediaServer.Close)

	dataURL := mediaServer.URL + "/rails/active_storage/blobs/redirect/media.jpg"
	attachmentJSON := fmt.Sprintf(`{
		"id": 77,
		"message_id": 101,
		"file_type": "image",
		"data_url": %q
	}`, dataURL)

	codechatSvc, requests := newCodechatServiceRecorder(t)
	relay := newRelayServiceForCodechatTest(codechatSvc)
	payload := mustDecodeChatwootWebhook(t, chatwootWebhookJSON("Imagem do local: **<"+mapsURL+">**", attachmentJSON))

	if err := relay.FromChatwoot(payload); err != nil {
		t.Fatalf("FromChatwoot returned error: %v", err)
	}

	if len(mediaRequests) != 1 || mediaRequests[0] != "HEAD /rails/active_storage/blobs/redirect/media.jpg" {
		t.Fatalf("expected one media HEAD probe, got %#v", mediaRequests)
	}
	req := requireSingleCodechatRequest(t, requests)
	if req.Path != "/message/sendMedia/test-instance" {
		t.Fatalf("expected sendMedia URL, got %q", req.Path)
	}

	var body struct {
		Number       string `json:"number"`
		MediaMessage struct {
			Mediatype string `json:"mediatype"`
			Caption   string `json:"caption"`
			Media     string `json:"media"`
		} `json:"mediaMessage"`
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		t.Fatalf("failed to decode CodeChat request body: %v", err)
	}
	if body.Number != "5549991158078" {
		t.Fatalf("expected normalized phone number, got %q", body.Number)
	}
	if body.MediaMessage.Mediatype != "image" {
		t.Fatalf("expected image mediatype, got %q", body.MediaMessage.Mediatype)
	}
	if body.MediaMessage.Media != dataURL {
		t.Fatalf("expected media URL %q, got %q", dataURL, body.MediaMessage.Media)
	}
	wantCaption := "Imagem do local: " + mapsURL
	if body.MediaMessage.Caption != wantCaption {
		t.Fatalf("expected normalized caption %q, got %q", wantCaption, body.MediaMessage.Caption)
	}
}

func TestRelayServiceFromChatwoot_ContractDocumentCaptionToCodechat(t *testing.T) {
	const mapsURL = "https://maps.app.goo.gl/UYZvebsBMEzk7cqk8"
	const documentURL = "https://chatwoot.test/rails/active_storage/blobs/redirect/contrato.pdf?disposition=attachment"

	attachmentJSON := fmt.Sprintf(`{
		"id": 78,
		"message_id": 102,
		"file_type": "file",
		"data_url": %q
	}`, documentURL)

	codechatSvc, requests := newCodechatServiceRecorder(t)
	relay := newRelayServiceForCodechatTest(codechatSvc)
	payload := mustDecodeChatwootWebhook(t, chatwootWebhookJSON("Documento: **<"+mapsURL+">**", attachmentJSON))

	if err := relay.FromChatwoot(payload); err != nil {
		t.Fatalf("FromChatwoot returned error: %v", err)
	}

	req := requireSingleCodechatRequest(t, requests)
	if req.Path != "/message/sendMedia/test-instance" {
		t.Fatalf("expected sendMedia URL, got %q", req.Path)
	}

	var body struct {
		Number       string `json:"number"`
		MediaMessage struct {
			Mediatype string `json:"mediatype"`
			FileName  string `json:"fileName"`
			Caption   string `json:"caption"`
			Media     string `json:"media"`
		} `json:"mediaMessage"`
	}
	if err := json.Unmarshal(req.Body, &body); err != nil {
		t.Fatalf("failed to decode CodeChat request body: %v", err)
	}
	if body.Number != "5549991158078" {
		t.Fatalf("expected normalized phone number, got %q", body.Number)
	}
	if body.MediaMessage.Mediatype != "document" {
		t.Fatalf("expected document mediatype, got %q", body.MediaMessage.Mediatype)
	}
	if body.MediaMessage.FileName != "contrato.pdf" {
		t.Fatalf("expected filename contrato.pdf, got %q", body.MediaMessage.FileName)
	}
	if body.MediaMessage.Media != documentURL {
		t.Fatalf("expected document URL %q, got %q", documentURL, body.MediaMessage.Media)
	}
	wantCaption := "Documento: " + mapsURL
	if body.MediaMessage.Caption != wantCaption {
		t.Fatalf("expected normalized caption %q, got %q", wantCaption, body.MediaMessage.Caption)
	}
}

func TestRelayServiceFromChatwoot_DoesNotCallCodechatForIgnoredWebhooks(t *testing.T) {
	const mapsURL = "https://maps.app.goo.gl/UYZvebsBMEzk7cqk8"

	tests := []struct {
		name   string
		mutate func(*dto.ChatwootWebhook)
	}{
		{
			name: "different event",
			mutate: func(payload *dto.ChatwootWebhook) {
				payload.Event = "conversation_updated"
			},
		},
		{
			name: "incoming message",
			mutate: func(payload *dto.ChatwootWebhook) {
				payload.MessageType = dto.Incoming
			},
		},
		{
			name: "private message",
			mutate: func(payload *dto.ChatwootWebhook) {
				payload.Private = true
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			codechatSvc, requests := newCodechatServiceRecorder(t)
			relay := newRelayServiceForCodechatTest(codechatSvc)
			payload := mustDecodeChatwootWebhook(t, chatwootWebhookJSON("Ignorar **<"+mapsURL+">**", ""))
			tc.mutate(&payload)

			if err := relay.FromChatwoot(payload); err != nil {
				t.Fatalf("FromChatwoot returned error: %v", err)
			}
			if len(*requests) != 0 {
				t.Fatalf("expected no CodeChat requests, got %#v", *requests)
			}
		})
	}
}

func sendChatwootTextThroughRelay(t *testing.T, content string, senderType string) string {
	t.Helper()

	var capturedText string
	codechatSvc := newCodechatServiceTestDouble(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/message/sendText/test-instance" {
			t.Fatalf("expected sendText URL, got %q", r.URL.Path)
		}

		var payload struct {
			TextMessage struct {
				Text string `json:"text"`
			} `json:"textMessage"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		capturedText = payload.TextMessage.Text

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	ctx := context.Background()
	relay := &RelayService{
		codechat: codechatSvc,
		ctx:      &ctx,
	}

	if err := relay.FromChatwoot(newChatwootTextWebhook(content, senderType)); err != nil {
		t.Fatalf("FromChatwoot returned error: %v", err)
	}

	return capturedText
}

func newChatwootTextWebhook(content string, senderType string) dto.ChatwootWebhook {
	return dto.ChatwootWebhook{
		Event:       "message_created",
		MessageType: dto.Outgoing,
		Private:     false,
		Conversation: dto.CWConversation{
			Meta: dto.CWMeta{
				Sender: dto.CWSenderMeta{
					Name:        "Contato Teste",
					PhoneNumber: "+554991158078",
				},
			},
			Messages: []dto.CWMessage{
				{Content: &content, SenderType: senderType},
			},
		},
	}
}

func newCodechatServiceTestDouble(t *testing.T, handler http.HandlerFunc) *CodechatService {
	t.Helper()

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	client, err := codechat.New(
		ts.URL,
		"global-token",
		codechat.WithInstanceToken("instance-token", "test-instance"),
	)
	if err != nil {
		t.Fatalf("failed to create codechat client: %v", err)
	}

	return &CodechatService{client: client}
}

type capturedCodechatRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   []byte
}

func newCodechatServiceRecorder(t *testing.T) (*CodechatService, *[]capturedCodechatRequest) {
	t.Helper()

	requests := []capturedCodechatRequest{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read CodeChat request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		requests = append(requests, capturedCodechatRequest{
			Method: r.Method,
			Path:   r.URL.Path,
			Header: r.Header.Clone(),
			Body:   body,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	t.Cleanup(ts.Close)

	client, err := codechat.New(
		ts.URL,
		"global-token",
		codechat.WithInstanceToken("instance-token", "test-instance"),
	)
	if err != nil {
		t.Fatalf("failed to create codechat client: %v", err)
	}

	return &CodechatService{client: client}, &requests
}

func newRelayServiceForCodechatTest(codechatSvc *CodechatService) *RelayService {
	ctx := context.Background()
	return &RelayService{
		codechat: codechatSvc,
		ctx:      &ctx,
	}
}

func requireSingleCodechatRequest(t *testing.T, requests *[]capturedCodechatRequest) capturedCodechatRequest {
	t.Helper()

	if len(*requests) != 1 {
		t.Fatalf("expected one CodeChat request, got %d: %#v", len(*requests), *requests)
	}
	return (*requests)[0]
}

func mustDecodeChatwootWebhook(t *testing.T, raw string) dto.ChatwootWebhook {
	t.Helper()

	var payload dto.ChatwootWebhook
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		t.Fatalf("failed to decode Chatwoot webhook fixture: %v\n%s", err, raw)
	}
	return payload
}

func chatwootWebhookJSON(content string, attachmentsJSON string) string {
	return fmt.Sprintf(`{
		"event": "message_created",
		"message_type": "outgoing",
		"private": false,
		"conversation": {
			"meta": {
				"sender": {
					"name": "Contato Teste",
					"phone_number": "+55 49 99115-8078"
				}
			},
			"messages": [
				{
					"id": 101,
					"content": %q,
					"content_type": "text",
					"message_type": 1,
					"sender_type": "AgentBot",
					"attachments": [%s]
				}
			]
		}
	}`, content, attachmentsJSON)
}
