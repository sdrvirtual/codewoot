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

func (c *Client) GetContactByPhone(ctx context.Context, phoneNumber string) (*dto.CWContact, error) {
	p := fmt.Sprintf("/api/v1/accounts/%d/contacts/search", c.accountId)

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

	// s, _ := json.MarshalIndent(raw, "", "\t")
	// fmt.Println(string(s))

	if len(out.Payload) > 0 {
		return &out.Payload[0], nil
	}
	return nil, nil
}

func (c *Client) CreateContact(ctx context.Context, params CreateContactParams) (*dto.CWContact, error) {
	if params.InboxID <= 0 {
		return nil, fmt.Errorf("inbox_id is required")
	}
	p := fmt.Sprintf("/api/v1/accounts/%d/contacts", c.accountId)
	req, err := c.newRequest(ctx, http.MethodPost, p, params)
	if err != nil {
		return nil, err
	}
	raw, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var out dto.CWContact
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode contact: %w", err)
	}
	return &out, nil
}

func (c *Client) CreateContactInbox(ctx context.Context, contactId int, params CreateContactInboxParams) (*dto.CWContactInbox, error) {
	p := fmt.Sprintf("/api/v1/accounts/%d/contacts/%d/contact_inboxes", c.accountId, contactId)
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
	p := fmt.Sprintf("/api/v1/accounts/%d/conversations", c.accountId)
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

// CreateMessage creates a message in a conversation.
// messageType should be "incoming" (from customer) or "outgoing" (from agent).
// Endpoint: POST /api/v1/accounts/{account_id}/conversations/{conversation_id}/messages
func (c *Client) CreateMessage(ctx context.Context, conversationID int, content, messageType string, private bool, contentAttributes map[string]any) (*dto.CWMessage, error) {
	if conversationID <= 0 {
		return nil, fmt.Errorf("conversation_id is required")
	}
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}
	if messageType == "" {
		messageType = "incoming"
	}
	body := map[string]any{
		"content":      content,
		"message_type": messageType,
		"private":      private,
	}
	if contentAttributes != nil {
		body["content_attributes"] = contentAttributes
	}

	p := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/messages", c.accountId, conversationID)
	req, err := c.newRequest(ctx, http.MethodPost, p, body)
	if err != nil {
		return nil, err
	}
	raw, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var msg dto.CWMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, fmt.Errorf("decode message: %w", err)
	}
	return &msg, nil
}

// SendAPIMessage is a convenience method that:
// 1) Creates a contact in the API inbox,
// 2) Starts a conversation using the returned source_id,
// 3) Sends an "incoming" message in that conversation.
func (c *Client) SendAPIMessage(ctx context.Context, inboxID int, name, phone, content string) (*dto.CWMessage, error) {
	contact, err := c.CreateContact(ctx, CreateContactParams{
		InboxID:     inboxID,
		Name:        name,
		PhoneNumber: phone,
	})
	if err != nil {
		return nil, err
	}
	if len(contact.ContactInboxes) == 0 || contact.ContactInboxes[0].SourceID == "" {
		return nil, fmt.Errorf("no source_id returned for contact")
	}
	sourceID := contact.ContactInboxes[0].SourceID

	convID, err := c.CreateConversation(ctx, sourceID, inboxID)
	if err != nil {
		return nil, err
	}

	return c.CreateMessage(ctx, convID, content, "incoming", false, nil)
}
