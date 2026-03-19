package audio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// TranscodeOggToMp3 re-encodes audio to OGG/Opus with WhatsApp-compatible
// parameters (mono, 48 kHz). The function name is historical; it actually
// produces OGG/Opus output — not MP3.
func TranscodeOggToMp3(oggfile io.Reader) (io.Reader, error) {
	var out bytes.Buffer

	err := ffmpeg.
		Input("pipe:0").
		Output("pipe:1", ffmpeg.KwArgs{
			"f":   "ogg",
			"c:a": "libopus",
			"b:a": "128k",
			"ac":  "1",
			"ar":  "48000",
		}).
		WithInput(oggfile).
		WithOutput(&out).
		Run()

	if err != nil {
		return nil, err
	}

	return &out, nil
}

// TranscodeToOgg converts any audio input to OGG/Opus with parameters that
// are compatible with WhatsApp voice notes on all devices including iPhone:
// mono channel, 48 kHz sample rate, Opus codec.
func TranscodeToOgg(audiofile io.Reader) (io.Reader, error) {
	var out bytes.Buffer

	err := ffmpeg.
		Input("pipe:0").
		Output("pipe:1", ffmpeg.KwArgs{
			"f":   "ogg",
			"c:a": "libopus",
			"b:a": "128k",
			"ac":  "1",
			"ar":  "48000",
		}).
		WithInput(audiofile).
		WithOutput(&out).
		Run()

	if err != nil {
		return nil, err
	}

	return &out, nil
}

func DownloadAndTranscodeToOgg(ctx context.Context, url string) (io.Reader, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download audio: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download audio: status %d", resp.StatusCode)
	}

	oggData, err := TranscodeToOgg(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to transcode audio: %w", err)
	}

	return oggData, nil
}
