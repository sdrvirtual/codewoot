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

func (c *Client) SendText(ctx context.Context, payload SendTextParams) (json.RawMessage, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendText/%s", url.PathEscape(c.instance))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	return jr, err
}

func (c *Client) SendWhatsappAudio(ctx context.Context, payload SendWhatsappAudioParams) (json.RawMessage, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	s, _ := json.MarshalIndent(payload, "", "\t")
	fmt.Println(string(s))
	p := fmt.Sprintf("/message/sendWhatsappAudio/%s", url.PathEscape(c.instance))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	return jr, err
}
