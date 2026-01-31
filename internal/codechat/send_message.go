package codechat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type CCAudioMessage struct {
	Audio string `json:"audio"`
}

type CCTextMessage struct {
	Text string `json:"text"`
}

type CCMediaMessage struct {
	Mediatype string `json:"mediatype,omitempty"`
	FileName  string `json:"fileName,omitempty"`
	Caption   string `json:"caption,omitempty"`
	Media     string `json:"media"`
}

type CCMessageOptions struct {
	ExternalAttributes string `json:"ExternalAttributes"`
	Delay              int    `json:"delay"`
	Presence           string `json:"presence"`
}

type SendTextParams struct {
	Number      string            `json:"number"`
	Options     *CCMessageOptions `json:"options,omitempty"`
	TextMessage CCTextMessage     `json:"textMessage"`
}

type SendWhatsappAudioParams struct {
	Number    string
	AudioFile io.Reader
	FileName  string
}

type SendMediaParams struct {
	Number       string            `json:"number"`
	Options      *CCMessageOptions `json:"options,omitempty"`
	MediaMessage CCMediaMessage    `json:"mediaMessage"`
}

func (c *Client) messageRequest(ctx context.Context, path string, payload any) (json.RawMessage, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/%s/%s", path, url.PathEscape(c.instance))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	return jr, err
}

func (c *Client) SendText(ctx context.Context, payload SendTextParams) (json.RawMessage, error) {
	return c.messageRequest(ctx, "sendText", payload)
}

func (c *Client) SendWhatsappAudio(ctx context.Context, payload SendWhatsappAudioParams) (json.RawMessage, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instanceName is required")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("number", payload.Number); err != nil {
		return nil, fmt.Errorf("failed to write number field: %w", err)
	}

	part, err := writer.CreateFormFile("attachment", payload.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, payload.AudioFile); err != nil {
		return nil, fmt.Errorf("failed to copy audio file: %w", err)
	}

	contentType := writer.FormDataContentType()
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	u := *c.baseURL
	u.Path = fmt.Sprintf("%s/message/sendWhatsappAudioFile/%s", c.baseURL.Path, url.PathEscape(c.instance))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("apikey", c.globalToken)
	if c.instanceToken != "" {
		req.Header.Set("Authorization", c.instanceToken)
	}

	jr, _, err := c.do(req)
	return jr, err
}

func (c *Client) SendMedia(ctx context.Context, payload SendMediaParams) (json.RawMessage, error) {
	return c.messageRequest(ctx, "sendMedia", payload)
}
