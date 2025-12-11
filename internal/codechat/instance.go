package codechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type CreateInstanceParams struct {
	InstanceName string `json:"instanceName"`
	Description  string `json:"description"`
}

type CreateInstanceResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Auth        struct {
		ID    int    `json:"id"`
		Token string `json:"token"`
	} `json:"Auth"`
}

type FetchInstanceResponse struct {
	ID                 int                `json:"id"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	ConnectionStatus   string             `json:"connectionStatus"`
	OwnerJid           string             `json:"ownerJid"`
	ProfilePicURL      string             `json:"profilePicUrl"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
	ExternalAttributes string             `json:"externalAttributes"`
	Webhook            SetWebhookResponse `json:"Webhook,omitempty"`
}

type ConnectInstanceResponse struct {
	Count  int    `json:"count"`
	Base64 string `json:"base64"`
}

func (c *Client) CreateInstance(ctx context.Context, payload CreateInstanceParams) (*CreateInstanceResponse, error) {
	p := "/instance/create"
	req, err := c.newRequest(ctx, http.MethodPost, p, payload)
	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var resp CreateInstanceResponse
	if err := json.Unmarshal(jr, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) FetchInstance(ctx context.Context) (*FetchInstanceResponse, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instance is required")
	}
	p := "/instance/fetchInstance/" + url.PathEscape(c.instance)
	req, err := c.newRequest(ctx, http.MethodGet, p, nil)
	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var resp FetchInstanceResponse
	if err := json.Unmarshal(jr, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ConnectInstance(ctx context.Context) (*ConnectInstanceResponse, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instance is required")
	}
	p := "/instance/connect/" + url.PathEscape(c.instance)
	req, err := c.newRequest(ctx, http.MethodGet, p, nil)
	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var resp ConnectInstanceResponse
	if err := json.Unmarshal(jr, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) LogoutInstance(ctx context.Context) (*json.RawMessage, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instance is required")
	}
	p := "/instance/logout/" + url.PathEscape(c.instance)
	req, err := c.newRequest(ctx, http.MethodDelete, p, nil)
	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	if err != nil {
		return nil, err
	}
	return &jr, nil
}

func (c *Client) DeleteInstance(ctx context.Context) (*json.RawMessage, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instance is required")
	}
	p := "/instance/delete/" + url.PathEscape(c.instance)
	req, err := c.newRequest(ctx, http.MethodDelete, p, nil)

	q := req.URL.Query()
	q.Set("force", "true")
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	if err != nil {
		return nil, err
	}
	return &jr, nil
}
