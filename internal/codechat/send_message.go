package codechat

import (
	"context"
	"encoding/json"
	"fmt"
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
	Number       string            `json:"number"`
	Options      *CCMessageOptions `json:"options,omitempty"`
	AudioMessage CCAudioMessage    `json:"audioMessage"`
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
	return c.messageRequest(ctx, "sendWhatsappAudio", payload)
}

func (c *Client) SendMedia(ctx context.Context, payload SendMediaParams) (json.RawMessage, error) {
	return c.messageRequest(ctx, "sendMedia", payload)
}
