package services

import (
	"context"
	"fmt"
	"log"

	"github.com/sdrvirtual/codewoot/internal/chatwoot"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/db"
	"github.com/sdrvirtual/codewoot/internal/domain"
	"github.com/sdrvirtual/codewoot/internal/dto"
)

type ChatwootService struct {
	cfg     *config.Config
	client  *chatwoot.Client
	inboxID int
}

type ConversationID int

func NewChatwootService(cfg *config.Config, session db.CodechatSession) *ChatwootService {
	// TODO: Colocar isso na db
	token := session.ChatwootToken
	accountID := int(session.ChatwootAccountID)
	inboxID := int(session.ChatwootInboxID)
	client, err := chatwoot.New(cfg.Chatwoot.URL, token, accountID)
	if err != nil {
		log.Fatal(err)
	}

	return &ChatwootService{cfg, client, inboxID}
}

func (c *ChatwootService) SetupContact(ctx context.Context, contact *domain.ContactInfo) (*dto.CWContact, error) {
	ctt, err := c.client.GetContactByPhone(ctx, contact.Phone)
	if err != nil {
		return nil, err
	}
	if ctt == nil {
		ctt, err = c.client.CreateContact(ctx, chatwoot.CreateContactParams{
			InboxID:     c.inboxID,
			Name:        contact.Name,
			PhoneNumber: contact.Phone},
		)
		if err != nil {
			return nil, err
		}
	}
	return ctt, nil
}

func (c *ChatwootService) setupInbox(ctx context.Context, contact *dto.CWContact) (*dto.CWContactInbox, error) {
	for _, ci := range contact.ContactInboxes {
		if ci.Inbox.ID == c.inboxID {
			return &ci, nil
		}
	}
	cttInbox, err := c.client.CreateContactInbox(ctx, contact.ID, chatwoot.CreateContactInboxParams{InboxID: c.inboxID})
	if err != nil {
		return nil, err
	}
	return cttInbox, nil
}

func (c *ChatwootService) setupConversation(ctx context.Context, contact *domain.ContactInfo) (ConversationID, error) {
	ctt, err := c.SetupContact(ctx, contact)
	if err != nil {
		return -1, err
	}
	if ctt == nil || ctt.ID < 1 {
		return  -1, fmt.Errorf("couldn't find or create account")
	}

	cttConv, err := c.client.GetContactConversations(ctx, ctt.ID)
	if err != nil {
		return -1, err
	}
	// Try to find an open conversation
	for _, conv := range cttConv {
		if conv.InboxID == c.inboxID {
			return ConversationID(conv.ID), nil
		}
	}

	// Conversation not found, or on another inbox
	cttInbox, err := c.setupInbox(ctx, ctt)
	if cttInbox.SourceID == "" {
		return -1, fmt.Errorf("source_id unavaliable")
	}
	if err != nil {
		return -1, err
	}
	convID, err := c.client.CreateConversation(ctx, cttInbox.SourceID, cttInbox.Inbox.ID)
	return ConversationID(convID), err
}

func (c *ChatwootService) SendMessage(ctx context.Context, contact domain.ContactInfo, message chatwoot.ChatwootClientMessage) error {
	id, err := c.setupConversation(ctx, &contact)
	message.ConversationID = int(id)
	if err != nil {
		return err
	}
	_, err = c.client.CreateMessage(ctx, message)
	return err
}
