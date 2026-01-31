package audio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func TranscodeOggToMp3(oggfile io.Reader) (io.Reader, error) {
	var out bytes.Buffer

	err := ffmpeg.
		Input("pipe:0").
		Output("pipe:1", ffmpeg.KwArgs{
			"f":   "ogg",
			"c:a": "libopus",
			"b:a": "128k",
		}).
		WithInput(oggfile).
		WithOutput(&out).
		Run()

	if err != nil {
		return nil, err
	}

	return &out, nil
}

func TranscodeToOgg(audiofile io.Reader) (io.Reader, error) {
	var out bytes.Buffer

	err := ffmpeg.
		Input("pipe:0").
		Output("pipe:1", ffmpeg.KwArgs{
			"f":   "ogg",
			"c:a": "libopus",
			"b:a": "128k",
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
