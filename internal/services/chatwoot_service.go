package services

import (
	"context"
	"fmt"
	"log"
	"log/slog"

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
	token := session.ChatwootToken
	accountID := int(session.ChatwootAccountID)
	inboxID := int(session.ChatwootInboxID)
	client, err := chatwoot.New(cfg.Chatwoot.URL, token, accountID)
	if err != nil {
		log.Fatal(err)
	}

	c := &ChatwootService{
		cfg:     cfg,
		client:  client,
		inboxID: inboxID,
	}

	return c
}

func (c *ChatwootService) SetupContact(ctx context.Context, contact *domain.ContactInfo) (*dto.CWContact, error) {
	slog.Info(fmt.Sprintf("Looking up contact with phone: %s", contact.Phone))
	ctt, err := c.client.GetContactByPhone(ctx, contact.Phone)
	if err != nil {
		return nil, err
	}
	if ctt == nil {
		slog.Info(fmt.Sprintf("Contact not found, creating new contact: %s (%s)", contact.Name, contact.Phone))
		ctt, err = c.client.CreateContact(ctx, chatwoot.CreateContactParams{
			InboxID:     c.inboxID,
			Name:        contact.Name,
			PhoneNumber: contact.Phone},
		)
		if err != nil {
			return nil, err
		}
		slog.Info(fmt.Sprintf("Contact created successfully with ID: %d", ctt.ID))
	} else {
		slog.Info(fmt.Sprintf("Found existing contact with ID: %d", ctt.ID))
	}
	return ctt, nil
}

func (c *ChatwootService) setupInbox(ctx context.Context, contact *dto.CWContact) (*dto.CWContactInbox, error) {
	slog.Info(fmt.Sprintf("Setting up inbox for contact ID: %d", contact.ID))
	for _, ci := range contact.ContactInboxes {
		if ci.Inbox.ID == c.inboxID {
			slog.Info(fmt.Sprintf("Found existing contact inbox with source ID: %s", ci.SourceID))
			return &ci, nil
		}
	}
	slog.Info(fmt.Sprintf("Contact inbox not found, creating new one for inbox ID: %d", c.inboxID))
	cttInbox, err := c.client.CreateContactInbox(ctx, contact.ID, chatwoot.CreateContactInboxParams{InboxID: c.inboxID})
	if err != nil {
		return nil, err
	}
	slog.Info(fmt.Sprintf("Contact inbox created successfully with source ID: %s", cttInbox.SourceID))
	return cttInbox, nil
}

func (c *ChatwootService) setupConversation(ctx context.Context, contact *domain.ContactInfo) (ConversationID, error) {
	slog.Info(fmt.Sprintf("Setting up conversation for contact: %s (%s)", contact.Name, contact.Phone))
	ctt, err := c.SetupContact(ctx, contact)
	if err != nil {
		return -1, err
	}
	if ctt == nil || ctt.ID < 1 {
		return -1, fmt.Errorf("couldn't find or create account")
	}

	slog.Info(fmt.Sprintf("Fetching existing conversations for contact ID: %d", ctt.ID))
	cttConv, err := c.client.GetContactConversations(ctx, ctt.ID)
	if err != nil {
		return -1, err
	}
	// Try to find an open conversation
	for _, conv := range cttConv {
		if conv.InboxID == c.inboxID {
			slog.Info(fmt.Sprintf("Found existing conversation with ID: %d", conv.ID))
			return ConversationID(conv.ID), nil
		}
	}

	slog.Info("No existing conversation found, creating new one")
	// Conversation not found, or on another inbox
	cttInbox, err := c.setupInbox(ctx, ctt)
	if err != nil {
		return -1, err
	}
	if cttInbox.SourceID == "" {
		return -1, fmt.Errorf("source_id unavaliable")
	}
	convID, err := c.client.CreateConversation(ctx, cttInbox.SourceID, cttInbox.Inbox.ID)
	if err != nil {
		return ConversationID(convID), err
	}
	slog.Info(fmt.Sprintf("Conversation created successfully with ID: %d", convID))
	return ConversationID(convID), err
}

func (c *ChatwootService) SendMessage(ctx context.Context, contact domain.ContactInfo, message chatwoot.ChatwootClientMessage) error {
	slog.Info(fmt.Sprintf("Sending message to contact: %s (%s)", contact.Name, contact.Phone))
	id, err := c.setupConversation(ctx, &contact)
	message.ConversationID = int(id)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Creating message in conversation ID: %d", id))
	msg, err := c.client.CreateMessage(ctx, message)
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Message created successfully with ID: %d", msg.ID))
	return nil
}
