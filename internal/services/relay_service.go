package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/sdrvirtual/codewoot/internal/chatwoot"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/domain"
	"github.com/sdrvirtual/codewoot/internal/dto"
)

type RelayService struct {
	cfg      *config.Config
	codechat *CodechatService
	chatwoot *ChatwootService
}

func NewRelayService(cfg *config.Config) *RelayService {
	// codechatInstanceToken := ""
	// instance := "codechat_v1"
	return &RelayService{
		cfg:      cfg,
		codechat: NewCodechatService(cfg),
		chatwoot: NewChatwootService(cfg),
	}
}

func (r *RelayService) FromCodechat(payload dto.CodechatWebhook) error {
	if payload.Data.IsGroup || payload.Data.KeyFromMe || payload.Event != "messages.upsert" {
		return nil
	}

	ctx := context.TODO()

	contact := domain.ContactInfo{
		Name:  payload.Data.PushName,
		Phone: "+" + strings.Split(payload.Data.KeyRemoteJid, "@")[0],
	}

	message := chatwoot.NewChatwootClientMessage()

	switch content := payload.Data.Content.(type) {
	case dto.CodechatTextContent:
		message.Text = content.Text
	case dto.CodechatAudioContent:
		audioData, err := r.codechat.GetAudioContent(ctx, &payload.Data)
		if err != nil {
			return err
		}
		message.FileType = "audio"
		message.Attachment = audioData
	case dto.CodechatImageContent:
		return fmt.Errorf("received message with images")
	}

	return r.chatwoot.SendMessage(ctx, contact, message)
}

func (r *RelayService) FromChatwoot(payload dto.ChatwootWebhook) error {
	if payload.Event != "message_created" || payload.MessageType != "outgoing" || payload.Private {
		return nil
	}

	contact := domain.ContactInfo{
		Name:  payload.Conversation.Meta.Sender.Name,
		Phone: strings.TrimPrefix(payload.Conversation.Meta.Sender.PhoneNumber, "+"),
	}

	ctx := context.TODO()

	for _, m := range payload.Conversation.Messages {
		message := NewCodechatClientMessage()

		if m.Content != nil {
			message.Text = *m.Content
		}

		for _, a := range m.Attachments {
			switch a.FileType {
			case "audio":
				message.AudioURL = a.DataURL
			}
		}

		if err := r.codechat.SendMessage(ctx, contact, message); err != nil {
			return err
		}
	}
	return nil
}
