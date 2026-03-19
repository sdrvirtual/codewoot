package codechat

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestSendWhatsappAudio_ContentType verifies that the multipart attachment
// part is sent with Content-Type "audio/ogg; codecs=opus" instead of the
// default "application/octet-stream" that Go's CreateFormFile would use.
// This is critical for WhatsApp on iPhone to correctly identify the audio
// as an OGG/Opus voice note and enable inline playback.
func TestSendWhatsappAudio_ContentType(t *testing.T) {
	audioData := []byte("fake-ogg-opus-data")

	var capturedBody []byte
	var capturedContentType string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		capturedBody = body
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	client, err := New(ts.URL, "test-global-token", WithInstanceToken("test-instance-token", "test-instance"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	params := SendWhatsappAudioParams{
		Number:    "5511999999999",
		AudioFile: bytes.NewReader(audioData),
		FileName:  "audio.ogg",
	}

	_, err = client.SendWhatsappAudio(t.Context(), params)
	if err != nil {
		t.Fatalf("SendWhatsappAudio returned error: %v", err)
	}

	// Parse the multipart body to inspect the attachment part's Content-Type
	mediaType, params2, err := mime.ParseMediaType(capturedContentType)
	if err != nil {
		t.Fatalf("failed to parse Content-Type header: %v", err)
	}
	if mediaType != "multipart/form-data" {
		t.Fatalf("expected multipart/form-data, got %s", mediaType)
	}

	boundary := params2["boundary"]
	reader := multipart.NewReader(bytes.NewReader(capturedBody), boundary)

	foundAttachment := false
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("error reading multipart part: %v", err)
		}

		if part.FormName() == "attachment" {
			foundAttachment = true
			partContentType := part.Header.Get("Content-Type")

			if !strings.Contains(partContentType, "audio/ogg") {
				t.Errorf("attachment Content-Type should contain 'audio/ogg', got: %s", partContentType)
			}
			if !strings.Contains(partContentType, "codecs=opus") {
				t.Errorf("attachment Content-Type should contain 'codecs=opus', got: %s", partContentType)
			}

			// Verify the body is correct
			body, _ := io.ReadAll(part)
			if !bytes.Equal(body, audioData) {
				t.Errorf("attachment body mismatch")
			}
		}
	}

	if !foundAttachment {
		t.Error("no attachment part found in multipart body")
	}
}

// TestSendWhatsappAudio_NumberField verifies that the phone number field
// is correctly included in the multipart form.
func TestSendWhatsappAudio_NumberField(t *testing.T) {
	var capturedBody []byte
	var capturedContentType string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContentType = r.Header.Get("Content-Type")
		body, _ := io.ReadAll(r.Body)
		capturedBody = body
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	client, err := New(ts.URL, "test-token", WithInstanceToken("inst-token", "inst"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	expectedNumber := "5511987654321"
	params := SendWhatsappAudioParams{
		Number:    expectedNumber,
		AudioFile: bytes.NewReader([]byte("audio")),
		FileName:  "voice.ogg",
	}

	_, err = client.SendWhatsappAudio(t.Context(), params)
	if err != nil {
		t.Fatalf("SendWhatsappAudio returned error: %v", err)
	}

	_, params2, _ := mime.ParseMediaType(capturedContentType)
	boundary := params2["boundary"]
	reader := multipart.NewReader(bytes.NewReader(capturedBody), boundary)

	foundNumber := false
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("error reading multipart part: %v", err)
		}

		if part.FormName() == "number" {
			body, _ := io.ReadAll(part)
			if string(body) != expectedNumber {
				t.Errorf("expected number %q, got %q", expectedNumber, string(body))
			}
			foundNumber = true
		}
	}

	if !foundNumber {
		t.Error("number field not found in multipart body")
	}
}

// TestSendWhatsappAudio_InstanceRequired verifies that an error is
// returned when instanceName is empty.
func TestSendWhatsappAudio_InstanceRequired(t *testing.T) {
	client, err := New("http://localhost:9999", "token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	params := SendWhatsappAudioParams{
		Number:    "5511999999999",
		AudioFile: bytes.NewReader([]byte("audio")),
		FileName:  "audio.ogg",
	}

	_, err = client.SendWhatsappAudio(t.Context(), params)
	if err == nil {
		t.Error("expected error when instance is empty")
	}
	if !strings.Contains(err.Error(), "instanceName is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestSendWhatsappAudio_URL verifies the request URL contains the correct
// endpoint and instance name.
func TestSendWhatsappAudio_URL(t *testing.T) {
	var capturedURL string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	instanceName := "my-whatsapp-instance"
	client, err := New(ts.URL, "token", WithInstanceToken("itok", instanceName))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	params := SendWhatsappAudioParams{
		Number:    "5511999999999",
		AudioFile: bytes.NewReader([]byte("audio")),
		FileName:  "audio.ogg",
	}

	_, err = client.SendWhatsappAudio(t.Context(), params)
	if err != nil {
		t.Fatalf("SendWhatsappAudio returned error: %v", err)
	}

	expected := "/message/sendWhatsappAudioFile/" + instanceName
	if capturedURL != expected {
		t.Errorf("expected URL path %q, got %q", expected, capturedURL)
	}
}
