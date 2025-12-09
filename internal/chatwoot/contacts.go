package chatwoot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sdrvirtual/codewoot/internal/dto"
)

func (c *Client) GetContactByPhone(ctx context.Context, phoneNumber string) (*dto.CWContact, error) {
	p := fmt.Sprintf("/api/v1/accounts/%d/contacts/search", c.accountID)

	req, err := c.newRequest(ctx, http.MethodGet, p, nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("q", phoneNumber)
	req.URL.RawQuery = q.Encode()

	raw, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var out struct {
		Payload []dto.CWContact `json:"payload"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode contact: %w", err)
	}

	if len(out.Payload) > 0 {
		return &out.Payload[0], nil
	}
	return nil, nil
}

func (c *Client) CreateContact(ctx context.Context, params CreateContactParams) (*dto.CWContact, error) {
	if params.InboxID <= 0 {
		return nil, fmt.Errorf("inbox_id is required")
	}
	p := fmt.Sprintf("/api/v1/accounts/%d/contacts", c.accountID)
	req, err := c.newRequest(ctx, http.MethodPost, p, params)
	if err != nil {
		return nil, err
	}
	raw, err := c.do(req)

	if err != nil {
		return nil, err
	}
	var out struct {
		Payload struct {
			Contact *dto.CWContact `json:"contact"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode contact: %w", err)
	}
	return out.Payload.Contact, nil
}

func (c *Client) CreateContactInbox(ctx context.Context, contactID int, params CreateContactInboxParams) (*dto.CWContactInbox, error) {
	p := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d/contact_inboxes", c.accountID, contactID)
	req, err := c.newRequest(ctx, http.MethodPost, p, params)
	if err != nil {
		return nil, err
	}
	raw, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var out dto.CWContactInbox
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode contact: %w", err)
	}

	return &out, nil
}
