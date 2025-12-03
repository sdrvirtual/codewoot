package chatwoot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sdrvirtual/codewoot/internal/dto"
)

type CreateContactParams struct {
	InboxID              int            `json:"inbox_id"` // required for API channel
	Name                 string         `json:"name,omitempty"`
	Email                string         `json:"email,omitempty"`
	PhoneNumber          string         `json:"phone_number,omitempty"`
	Thumbnail            string         `json:"thumbnail,omitempty"`
	AdditionalAttributes map[string]any `json:"additional_attributes,omitempty"`
}

type CreateContactInboxParams struct {
	InboxID  int     `json:"inbox_id"`
	SourceID *string `json:"source_id"`
}

type ConversationRef struct {
	ID int `json:"id"`
}

func (c *Client) GetContactConversations(ctx context.Context, contactID int) ([]dto.CWConversation, error) {
	p := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d/conversations", c.accountID, contactID)
	req, err := c.newRequest(ctx, http.MethodGet, p, nil)
	if err != nil {
		return nil, err
	}
	raw, err := c.do(req)
	if err != nil {
		return nil, err
	}

	var out struct {
		Payload []dto.CWConversation `json:"payload"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode conversation: %w", err)
	}

	return out.Payload, nil
}

func (c *Client) CreateConversation(ctx context.Context, sourceID string, inboxID int) (int, error) {
	if sourceID == "" {
		return 0, fmt.Errorf("source_id is required")
	}
	if inboxID <= 0 {
		return 0, fmt.Errorf("inbox_id is required")
	}
	body := map[string]any{
		"source_id": sourceID,
		"inbox_id":  inboxID,
	}
	p := fmt.Sprintf("/api/v1/accounts/%d/conversations", c.accountID)
	req, err := c.newRequest(ctx, http.MethodPost, p, body)
	if err != nil {
		return 0, err
	}
	raw, err := c.do(req)
	if err != nil {
		return 0, err
	}
	var ref ConversationRef
	if err := json.Unmarshal(raw, &ref); err != nil {
		return 0, fmt.Errorf("decode conversation: %w", err)
	}
	return ref.ID, nil
}
