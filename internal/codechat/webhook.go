package codechat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SetWebhookParams struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
	Events  struct {
		QrcodeUpdated             bool `json:"qrcodeUpdated"`
		MessagesSet               bool `json:"messagesSet"`
		MessagesUpsert            bool `json:"messagesUpsert"`
		MessagesUpdated           bool `json:"messagesUpdated"`
		SendMessage               bool `json:"sendMessage"`
		ContactsSet               bool `json:"contactsSet"`
		ContactsUpsert            bool `json:"contactsUpsert"`
		ContactsUpdated           bool `json:"contactsUpdated"`
		ChatsSet                  bool `json:"chatsSet"`
		ChatsUpsert               bool `json:"chatsUpsert"`
		ChatsUpdated              bool `json:"chatsUpdated"`
		ChatsDeleted              bool `json:"chatsDeleted"`
		PresenceUpdated           bool `json:"presenceUpdated"`
		GroupsUpsert              bool `json:"groupsUpsert"`
		GroupsUpdated             bool `json:"groupsUpdated"`
		GroupsParticipantsUpdated bool `json:"groupsParticipantsUpdated"`
		ConnectionUpdated         bool `json:"connectionUpdated"`
		StatusInstance            bool `json:"statusInstance"`
		RefreshToken              bool `json:"refreshToken"`
	} `json:"events"`
}

type SetWebhookResponse struct {
	SetWebhookParams
	ID int `json:"id"`
}

func NewSetWebhookParams(url string) *SetWebhookParams {
	p := &SetWebhookParams{
		Enabled: true,
		URL:     url,
	}
	p.Events.MessagesUpsert = true
	return p
}

func (c *Client) SetWebhook(ctx context.Context, params *SetWebhookParams) (*SetWebhookResponse, error) {
	if c.instance == "" {
		return nil, fmt.Errorf("instance is required")
	}
	p := "/webhook/set/" + url.PathEscape(c.instance)
	req, err := c.newRequest(ctx, http.MethodPut, p, params)
	if err != nil {
		return nil, err
	}
	jr, _, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var resp SetWebhookResponse
	if err := json.Unmarshal(jr, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
