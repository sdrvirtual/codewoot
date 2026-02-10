package services

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdrvirtual/codewoot/internal/chatwoot"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/db"
	"github.com/sdrvirtual/codewoot/internal/domain"
	"github.com/sdrvirtual/codewoot/internal/dto"
	"github.com/sdrvirtual/codewoot/internal/utils"
)

var (
	globalPhoneMutexes = make(map[string]*sync.Mutex)
	globalMapMutex     sync.Mutex
)

type RelayService struct {
	cfg      *config.Config
	codechat *CodechatService
	chatwoot *ChatwootService
	ctx      *context.Context
}

func NewRelayService(ctx context.Context, cfg *config.Config, p *pgxpool.Pool, session string) (*RelayService, error) {
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
		ctx:      &ctx,
	}, nil
}

func (r *RelayService) getPhoneMutex(phone string) *sync.Mutex {
	globalMapMutex.Lock()
	defer globalMapMutex.Unlock()

	if _, exists := globalPhoneMutexes[phone]; !exists {
		globalPhoneMutexes[phone] = &sync.Mutex{}
	}

	return globalPhoneMutexes[phone]
}

func (r *RelayService) FromCodechat(payload dto.CodechatWebhook) error {
	if payload.Data.IsGroup || payload.Data.KeyFromMe || payload.Event != "messages.upsert" {
		return nil
	}

	if payload.Data.KeyRemoteJid == payload.Instance.OwnerJid {
		// This should not happen, but it does ;(
		return nil
	}

	// TODO: Handle deleting messages

	phone, err := utils.ValidatePhone(strings.Split(payload.Data.KeyRemoteJid, "@")[0])
	if err != nil {
		return err
	}

	phoneMutex := r.getPhoneMutex(phone)
	phoneMutex.Lock()
	defer phoneMutex.Unlock()

	contact := domain.ContactInfo{
		Name:  payload.Data.PushName,
		Phone: "+" + phone,
	}

	message := chatwoot.NewChatwootClientMessage()

	switch content := payload.Data.Content.(type) {
	case dto.CodechatTextContent:
		message.Text = content.Text
	case dto.CodechatAudioContent:
		audioData, err := r.codechat.GetAudioContent(*r.ctx, &payload.Data)
		if err != nil {
			return err
		}
		message.FileType = "audio"
		message.Attachment = audioData
	case dto.CodechatImageContent:
		imageData, err := r.codechat.GetMediaContent(*r.ctx, &payload.Data)
		if err != nil {
			return err
		}
		message.Text = content.Caption
		message.Attachment = imageData
	case dto.CodechatDocumentContent:
		documentData, err := r.codechat.GetMediaContent(*r.ctx, &payload.Data)
		if err != nil {
			return err
		}
		message.Text = content.Caption
		message.Attachment = documentData
	}

	return r.chatwoot.SendMessage(*r.ctx, contact, message)
}

func (r *RelayService) FromChatwoot(payload dto.ChatwootWebhook) error {
	if payload.Event != "message_created" || payload.MessageType != "outgoing" || payload.Private {
		return nil
	}

	phone, err := utils.ValidatePhone(strings.TrimPrefix(payload.Conversation.Meta.Sender.PhoneNumber, "+"))
	if err != nil {
		return err
	}

	phoneMutex := r.getPhoneMutex(phone)
	phoneMutex.Lock()
	defer phoneMutex.Unlock()

	contact := domain.ContactInfo{
		Name:  payload.Conversation.Meta.Sender.Name,
		Phone: phone,
	}

	// TODO: Handle deleting messages

	for _, m := range payload.Conversation.Messages {
		message := NewCodechatClientMessage()

		if m.Content != nil {
			message.Text = *m.Content
		}

		for _, a := range m.Attachments {
			switch a.FileType {
			case "audio":
				oggFile, err := r.codechat.TranscodeAudioFromURL(*r.ctx, *a.DataURL)
				if err != nil {
					return fmt.Errorf("failed to transcode audio: %w", err)
				}
				fileName := "audio.ogg"
				message.AudioFile = oggFile
				message.AudioFileName = &fileName
			case "image":
				message.MediaURL = a.DataURL
			case "file":
				message.FileURL = a.DataURL
				u, err := url.Parse(*a.DataURL)
				if err != nil {
					return err
				}
				filename := path.Base(u.Path)
				message.AttachmentName = &filename
			}
		}

		if err := r.codechat.SendMessage(*r.ctx, contact, message); err != nil {
			return err
		}
	}
	return nil
}
