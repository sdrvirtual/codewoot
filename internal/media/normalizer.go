package media

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type PreparedMedia struct {
	UseURL   bool
	File     io.Reader
	FileName string
	MimeType string
}

func PrepareImageFromURL(ctx context.Context, rawURL string) (*PreparedMedia, error) {
	mimeType, err := probeContentType(ctx, rawURL)
	if err != nil {
		return nil, err
	}

	switch mimeType {
	case "image/jpeg", "image/png":
		return &PreparedMedia{UseURL: true, MimeType: mimeType}, nil
	}

	if mimeType != "" && mimeType != "application/octet-stream" && !strings.HasPrefix(mimeType, "image/") {
		return nil, fmt.Errorf("expected image content type, got %s", mimeType)
	}

	data, err := download(ctx, rawURL)
	if err != nil {
		return nil, err
	}

	converted, err := transcodeImageToJPEG(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to transcode image to jpeg: %w", err)
	}

	return &PreparedMedia{
		File:     converted,
		FileName: "media.jpg",
		MimeType: "image/jpeg",
	}, nil
}

func PrepareVideoFromURL(ctx context.Context, rawURL string) (*PreparedMedia, error) {
	mimeType, err := probeContentType(ctx, rawURL)
	if err != nil {
		return nil, err
	}

	if mimeType == "video/mp4" {
		return &PreparedMedia{UseURL: true, MimeType: mimeType}, nil
	}

	if mimeType != "" && mimeType != "application/octet-stream" && !strings.HasPrefix(mimeType, "video/") {
		return nil, fmt.Errorf("expected video content type, got %s", mimeType)
	}

	data, err := download(ctx, rawURL)
	if err != nil {
		return nil, err
	}

	converted, err := transcodeVideoToMP4(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to transcode video to mp4: %w", err)
	}

	return &PreparedMedia{
		File:     converted,
		FileName: "media.mp4",
		MimeType: "video/mp4",
	}, nil
}

func probeContentType(ctx context.Context, rawURL string) (string, error) {
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, rawURL, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create media probe request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			if resp.Body != nil {
				resp.Body.Close()
			}
			if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
				return normalizeContentType(resp.Header.Get("Content-Type")), nil
			}
			lastErr = fmt.Errorf("media probe returned status %d", resp.StatusCode)
		} else {
			lastErr = err
		}

		if err := sleepBeforeRetry(ctx); err != nil {
			return "", err
		}
	}

	data, err := downloadRange(ctx, rawURL)
	if err != nil {
		if lastErr != nil {
			return "", fmt.Errorf("failed to probe media after retries: %w; fallback failed: %w", lastErr, err)
		}
		return "", err
	}

	return normalizeContentType(http.DetectContentType(data)), nil
}

func normalizeContentType(contentType string) string {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return strings.ToLower(strings.TrimSpace(contentType))
	}
	return strings.ToLower(mediaType)
}

func download(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create media download request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to download media: %w", err)
		} else {
			data, readErr := readSuccessfulResponse(resp)
			if readErr == nil {
				return data, nil
			}
			lastErr = readErr
		}

		if err := sleepBeforeRetry(ctx); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("failed to download media after retries: %w", lastErr)
}

func downloadRange(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create media probe fallback request: %w", err)
		}
		req.Header.Set("Range", "bytes=0-511")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to probe media: %w", err)
		} else {
			data, readErr := readSuccessfulResponse(resp)
			if readErr == nil {
				if len(data) > 512 {
					return data[:512], nil
				}
				return data, nil
			}
			lastErr = readErr
		}

		if err := sleepBeforeRetry(ctx); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("failed to probe media after retries: %w", lastErr)
}

func readSuccessfulResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("media request returned status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read media response: %w", err)
	}
	return data, nil
}

func sleepBeforeRetry(ctx context.Context) error {
	timer := time.NewTimer(500 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func transcodeImageToJPEG(input io.Reader) (io.Reader, error) {
	var out bytes.Buffer
	err := ffmpeg.
		Input("pipe:0").
		Output("pipe:1", ffmpeg.KwArgs{
			"f":       "image2",
			"vframes": "1",
			"vcodec":  "mjpeg",
		}).
		WithInput(input).
		WithOutput(&out).
		Run()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(out.Bytes()), nil
}

func transcodeVideoToMP4(input io.Reader) (io.Reader, error) {
	// WhatsApp mobile is stricter than the upload API: MP4 containers with odd
	// dimensions, high H.264 levels or unusual frame rates can be accepted by the
	// server but fail to download/play on the client. Normalize to a conservative
	// H.264 baseline profile at 30fps while padding odd dimensions by at most one
	// pixel instead of scaling/cropping the original content. Write to a temporary
	// file so ffmpeg can produce a normal MP4 with faststart metadata instead of a
	// fragmented pipe MP4, which WhatsApp may accept but fail to play on mobile.
	tmp, err := os.CreateTemp("", "codewoot-video-*.mp4")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary video file: %w", err)
	}
	tmpPath := tmp.Name()
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return nil, fmt.Errorf("failed to close temporary video file: %w", err)
	}
	if err := os.Remove(tmpPath); err != nil {
		return nil, fmt.Errorf("failed to prepare temporary video path: %w", err)
	}

	err = ffmpeg.
		Input("pipe:0").
		Output(tmpPath, ffmpeg.KwArgs{
			"f":         "mp4",
			"c:v":       "libx264",
			"profile:v": "baseline",
			"level":     "3.1",
			"c:a":       "aac",
			"pix_fmt":   "yuv420p",
			"preset":    "veryfast",
			"movflags":  "+faststart",
			"vf":        "fps=30,pad=ceil(iw/2)*2:ceil(ih/2)*2",
		}).
		WithInput(input).
		Run()
	if err != nil {
		_ = os.Remove(tmpPath)
		return nil, err
	}

	file, err := os.Open(tmpPath)
	if err != nil {
		_ = os.Remove(tmpPath)
		return nil, fmt.Errorf("failed to open transcoded video: %w", err)
	}

	return &cleanupReadCloser{File: file, path: tmpPath}, nil
}

type cleanupReadCloser struct {
	*os.File
	path string
}

func (c *cleanupReadCloser) Close() error {
	err := c.File.Close()
	removeErr := os.Remove(c.path)
	if err != nil {
		return err
	}
	return removeErr
}
