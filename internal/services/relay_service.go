package services

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/dto"
)

type RelayService struct {
	cfg *config.Config
	codechat *CodechatService
	chatwoot *ChatwootService
}

func NewRelayService(cfg *config.Config) *RelayService {
	// codechatInstanceToken := ""
	// instance := "codechat_v1"
	return &RelayService{
		cfg: cfg,
		codechat: NewCodechatService(cfg),
		chatwoot: NewChatwootService(cfg),
	}
}

func (r *RelayService) FromCodechat(payload dto.CodechatWebhook) {
	s, _ := json.MarshalIndent(payload, "", "\t")

	if payload.Data.IsGroup || payload.Data.KeyFromMe || payload.Event != "messages.upsert" {
		return
	}

	log.Println(string(s))
}

func (r *RelayService) FromChatwoot(payload dto.ChatwootWebhook) error {
	// s, _ := json.MarshalIndent(payload, "", "\t")
	// fmt.Println(string(s))
	if payload.Event != "message_created" || payload.MessageType != "outgoing" || payload.Private {
		return nil
	}

	number := strings.TrimPrefix(payload.Conversation.Meta.Sender.PhoneNumber, "+")

	for _, m := range payload.Conversation.Messages {
		err := r.codechat.SendTextMessage(number, m.Content)
		if err != nil {
			return err
		}
	}
	return nil
}
