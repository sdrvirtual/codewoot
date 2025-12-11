package services

import (
	"context"
	"fmt"
	"net/url"
	"path"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdrvirtual/codewoot/internal/codechat"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/db"
)

type SessionService struct {
	cfg    *config.Config
	client *codechat.Client
	ctx    *context.Context
	db     *db.Queries
}

type Option func(*SessionService) error

func WithInstance(instanceToken string, instance string) Option {
	return func(c *SessionService) error {
		client, err := codechat.New(
			c.cfg.Codechat.URL,
			c.cfg.Codechat.GlobalToken,
			codechat.WithInstanceToken(instanceToken, instance),
		)
		c.client = client
		return err
	}
}

func NewSessionService(ctx context.Context, cfg *config.Config, p *pgxpool.Pool, opts ...Option) (*SessionService, error) {
	client, err := codechat.New(cfg.Codechat.URL, cfg.Codechat.GlobalToken)
	if err != nil {
		return nil, err
	}
	db := db.New(p)
	s := &SessionService{
		cfg:    cfg,
		client: client,
		ctx:    &ctx,
		db:     db,
	}

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (s *SessionService) CreateSession(instanceID *string, description *string, token string, inboxID, accountID int) (*db.CodechatSession, error) {
	var sessionUUID pgtype.UUID
	var err error

	if instanceID != nil {
		err = sessionUUID.Scan(*instanceID)
	} else {
		u, _ := uuid.NewUUID()
		err = sessionUUID.Scan(u.String())
	}
	if err != nil {
		return nil, err
	}

	_, err = s.db.GetSessionBySessionId(*s.ctx, sessionUUID)
	if err == nil {
		return nil, fmt.Errorf("session exists")
	}
	if err.Error() != "no rows in result set" {
		return nil, err
	}
	cip := &codechat.CreateInstanceParams{
		InstanceName: sessionUUID.String(),
		Description:  *description,
	}
	instance, err := s.client.CreateInstance(*s.ctx, *cip)
	if err != nil {
		return nil, err
	}
	session, err := s.db.CreateSession(*s.ctx, db.CreateSessionParams{
		SessionID:              sessionUUID,
		ChatwootToken:          token,
		ChatwootInboxID:        int32(inboxID),
		ChatwootAccountID:      int32(accountID),
		CodechatInstance:       instance.Name,
		CodechatInstcanceToken: instance.Auth.Token,
	})
	if err != nil {
		return nil, err
	}
	return &session, err
}

func (s *SessionService) ConnectSession() (*string, error) {
	i, err := s.client.ConnectInstance(*s.ctx)
	if err != nil {
		return nil, err
	}
	return &i.Base64, nil
}

func (s *SessionService) SetWebhook(instance string) error {
	u, err := url.Parse(s.cfg.Server.URL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "/codechat/webhook/", instance)
	_, err = s.client.SetWebhook(*s.ctx, codechat.NewSetWebhookParams(u.String()))
	return err
}

func (s *SessionService) FetchInstance() (*codechat.FetchInstanceResponse, error) {
	r, err := s.client.FetchInstance(*s.ctx)
	return r, err
}

func (s *SessionService) DeleteInstance() error {
	r, err := s.client.LogoutInstance(*s.ctx)
	if err != nil {
		return err
	}
	fmt.Printf("r: %v\n", r)
	r, err = s.client.DeleteInstance(*s.ctx)
	fmt.Printf("r: %v\n", r)
	return err
}
