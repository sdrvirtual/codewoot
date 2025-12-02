package codechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) SendText(ctx context.Context, payload any) (json.RawMessage, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendText/%s", url.PathEscape(c.instance))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *Client) SendMedia(ctx context.Context, instanceName string, payload any) (json.RawMessage, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendMedia/%s", url.PathEscape(instanceName))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *Client) SendMediaFile(ctx context.Context, instanceName string, payload any) (json.RawMessage, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendMediaFile/%s", url.PathEscape(instanceName))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *Client) SendLocation(ctx context.Context, instanceName string, payload any) (json.RawMessage, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendLocation/%s", url.PathEscape(instanceName))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *Client) SendContact(ctx context.Context, instanceName string, payload any) (json.RawMessage, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendContact/%s", url.PathEscape(instanceName))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *Client) SendReaction(ctx context.Context, instanceName string, payload any) (json.RawMessage, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendReaction/%s", url.PathEscape(instanceName))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *Client) SendWhatsAppAudio(ctx context.Context, instanceName string, payload any) (json.RawMessage, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendWhatsAppAudio/%s", url.PathEscape(instanceName))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *Client) SendWhatsAppAudioFile(ctx context.Context, instanceName string, payload any) (json.RawMessage, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendWhatsAppAudioFile/%s", url.PathEscape(instanceName))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *Client) SendButtons(ctx context.Context, instanceName string, payload any) (json.RawMessage, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendButtons/%s", url.PathEscape(instanceName))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}

func (c *Client) SendList(ctx context.Context, instanceName string, payload any) (json.RawMessage, error) {
	if instanceName == "" {
		return nil, fmt.Errorf("instanceName is required")
	}
	p := fmt.Sprintf("/message/sendList/%s", url.PathEscape(instanceName))
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	return c.doJSON(req)
}
