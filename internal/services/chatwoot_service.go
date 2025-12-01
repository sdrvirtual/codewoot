package services

import (
	"context"
	"log"

	"github.com/sdrvirtual/codewoot/internal/chatwoot"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/dto"
)

type ChatwootService struct {
	cfg *config.Config
	client *chatwoot.Client
	inboxId int
}

func NewChatwootService(cfg *config.Config) *ChatwootService {
	// TODO: Colocar isso na db
	token := "7ufe5YTVz6gDrCYfVRV18uVr"
	account_id := 2
	inboxID := 1
	client, err := chatwoot.New(cfg.Chatwoot.URL, token, account_id)
	if err != nil {
		log.Fatal(err)
	}

	return &ChatwootService{cfg, client, inboxID}
}

func (c *ChatwootService) setupInbox(phoneNumber string) (*dto.CWContactInbox, error) {
	ctx := context.TODO()
	ctt, err := c.client.GetContactByPhone(ctx, phoneNumber)
	if err != nil {
		log.Fatalln(err)
	}
	for _, ci := range ctt.ContactInboxes {
		if ci.Inbox.ID == c.inboxId {
			return &ci, nil
		}
	}
	cttInbox, err := c.client.CreateContactInbox(ctx, ctt.ID, chatwoot.CreateContactInboxParams{InboxID: c.inboxId})
	if err != nil {
		return nil, err
	}
	return cttInbox, nil
}

func (c *ChatwootService) SendTextMessage(toNumber string, text string) {
	
}

