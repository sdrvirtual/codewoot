package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sdrvirtual/codewoot/internal/chatwoot"
	"github.com/sdrvirtual/codewoot/internal/config"
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
	name := payload.Data.PushName
	phone := "+" + strings.Split(payload.Data.KeyRemoteJid, "@")[0]

	contact := ContactInfo{
		Name:  name,
		Phone: phone,
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

	contact := ContactInfo{
		Name:  payload.Conversation.Meta.Sender.Name,
		Phone: strings.TrimPrefix(payload.Conversation.Meta.Sender.PhoneNumber, "+"),
	}

	s, _ := json.MarshalIndent(payload, "", "\t")
	fmt.Println(string(s))

	for _, m := range payload.Conversation.Messages {
		// err := r.codechat.SendTextMessage(number, m.Content)
		// if err != nil {
		// 	return err
		// }
	}
	return nil
}
