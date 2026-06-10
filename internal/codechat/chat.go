package codechat

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"

	"github.com/sdrvirtual/codewoot/internal/dto"
)

type WhatsappNumbersParams struct {
	Numbers []string `json:"numbers"`
}

type OnWhatsAppResponse struct {
	JID    string `json:"jid"`
	Exists bool   `json:"exists"`
	LID    string `json:"lid"`
	Name   string `json:"name"`
}

func (c *Client) WhatsappNumbers(ctx context.Context, numbers []string) ([]OnWhatsAppResponse, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/chat/whatsappNumbers/%s", url.PathEscape(c.instance))
	req, err := c.newRequest(ctx, http.MethodPost, p, WhatsappNumbersParams{Numbers: numbers})
	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var resp []OnWhatsAppResponse
	if err := json.Unmarshal(jr, &resp); err != nil {
		return nil, fmt.Errorf("decode whatsapp numbers: %w", err)
	}
	return resp, nil
}

func (c *Client) GetMediaData(ctx context.Context, message *dto.CodechatData) (*dto.FileData, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/chat/mediaData/%s", url.PathEscape(c.instance))
	req, err := c.newRequest(ctx, http.MethodPost, p, map[string]string{"keyId": message.KeyID})

	q := req.URL.Query()
	q.Set("binary", "true")
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}
	_, r, err := c.do(req)
	if err != nil {
		return nil, err
	}

	_, params, err := mime.ParseMediaType(r.Header.Get("Content-Disposition"))
	if err != nil {
		return nil, err
	}

	fileData := dto.FileData{
		Name:     params["filename"],
		Mimetype: r.Header.Get("Content-Type"),
		File:     r.Body,
	}

	return &fileData, nil
}
