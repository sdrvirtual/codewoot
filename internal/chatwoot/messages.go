package chatwoot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/sdrvirtual/codewoot/internal/dto"
)

type ChatwootClientMessage struct {
	Text           string
	ConversationID int
	MessageType    dto.CWMessageType
	Private        bool
	Attachment     *dto.FileData
}

func NewChatwootClientMessage() ChatwootClientMessage {
	return ChatwootClientMessage{
		Private: false,
		MessageType: dto.Incoming,
	}
}

func (c *Client) CreateMessage(ctx context.Context, message ChatwootClientMessage) (*dto.CWMessage, error) {
	p := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/messages", c.accountId, message.ConversationID)

	// params := message.toAPIMessage()
	var buf bytes.Buffer

	mw := multipart.NewWriter(&buf)

	if message.Text != "" {
		fw, err := mw.CreateFormField("content")
		if err != nil {
			return nil, err
		}
		_, err = fw.Write([]byte(message.Text))
		if err != nil {
			return nil, err
		}
	}
	if message.MessageType != "" {
		fw, err := mw.CreateFormField("message_type")
		if err != nil {
			return nil, err
		}
		_, err = fw.Write([]byte(message.MessageType))
		if err != nil {
			return nil, err
		}
	}
	if message.Attachment != nil {
		fw, err := mw.CreateFormFile("attachments[]", message.Attachment.Name)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(fw, message.Attachment.File)
		if err != nil {
			return nil, err
		}
	}

	mw.Close()

	req, err := c.newRequest(ctx, http.MethodPost, p, io.Reader(&buf))

	req.Header.Set("Content-Type", mw.FormDataContentType())

	if err != nil {
		return nil, err
	}
	raw, err := c.do(req)
	if err != nil {
		return nil, err
	}
	var out dto.CWMessage
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}
	return &out, nil
}
