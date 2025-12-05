package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdrvirtual/codewoot/internal/chatwoot"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/db"
	"github.com/sdrvirtual/codewoot/internal/domain"
	"github.com/sdrvirtual/codewoot/internal/dto"
)

type RelayService struct {
	cfg      *config.Config
	codechat *CodechatService
	chatwoot *ChatwootService
}

func NewRelayService(ctx context.Context,cfg *config.Config, p *pgxpool.Pool, session string) (*RelayService, error) {
	q := db.New(p)
	var sessionUUID pgtype.UUID
	err := sessionUUID.Scan(session)
	if err != nil {
		return nil, err
	}
	sessionObj, err := q.GetSessionBySessionId(ctx, sessionUUID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("session %s does not exist", session)
		}
		return nil, err
	}

	return &RelayService{
		cfg:      cfg,
		codechat: NewCodechatService(cfg, sessionObj),
		chatwoot: NewChatwootService(cfg, sessionObj),
	}, nil
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
			fmt.Println(a.FileType)
			switch a.FileType {
			case "audio":
				message.AudioURL = a.DataURL
			case "image":
				message.MediaURL = a.DataURL
			}
		}

		if err := r.codechat.SendMessage(ctx, contact, message); err != nil {
			return err
		}
	}
	return nil
}
