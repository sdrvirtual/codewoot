package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sdrvirtual/codewoot/internal/config"
	"github.com/sdrvirtual/codewoot/internal/db"
	"github.com/sdrvirtual/codewoot/internal/dto"
	"github.com/sdrvirtual/codewoot/internal/services"
)

type CreateSessionResponse struct {
	ID                   int    `json:"id"`
	SessionID            string `json:"session_id"`
	ChatwootInboxWebhook string `json:"chatwoot_inbox_webhook"`
}

func (rd *CreateSessionResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func newCreateSessionResponse(cfg *config.Config, session db.CodechatSession) *CreateSessionResponse {
	u, _ := url.Parse(cfg.Server.URL)
	u.Path = path.Join(u.Path, "/chatwoot/webhook/", session.SessionID.String())
	return &CreateSessionResponse{
		ID:                   int(session.ID),
		SessionID:            session.SessionID.String(),
		ChatwootInboxWebhook: u.String(),
	}
}

type StatusSessionResponse struct {
	CreateSessionResponse
	Status string `json:"status"`
}

func (rd *StatusSessionResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func newStatusSessionResponse(cfg *config.Config, session db.CodechatSession, status string) *StatusSessionResponse {
	return &StatusSessionResponse{
		CreateSessionResponse: *newCreateSessionResponse(cfg, session),
		Status:                status,
	}
}

type ConnectSessionResponse struct {
	Base64 string `json:"base64"`
}

func (rd *ConnectSessionResponse) Render(w http.ResponseWriter, r *http.Request) error { return nil }

func newConnectSessionResponse(base64 string) *ConnectSessionResponse {
	return &ConnectSessionResponse{
		Base64: base64,
	}
}

func CreateSession(cfg *config.Config, p *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2MB

		var payload dto.CreateSession

		dec := json.NewDecoder(r.Body)

		if err := dec.Decode(&payload); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("invalid payload", err.Error()))
			return
		}

		sessionService, err := services.NewSessionService(r.Context(), cfg, p)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("error creating service", err.Error()))
			return
		}

		session, err := sessionService.CreateSession(
			payload.SessionID,
			payload.Description,
			payload.Chatwoot.Token,
			payload.Chatwoot.InboxID,
			payload.Chatwoot.AccountID,
		)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("error creating session", err.Error()))
			return
		}

		sessionService, err = services.NewSessionService(
			r.Context(),
			cfg,
			p,
			services.WithInstance(session.CodechatInstcanceToken, session.CodechatInstance),
		)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("error creating service", err.Error()))
			return
		}
		if err = sessionService.SetWebhook(session.SessionID.String()); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("error configuring webhook", err.Error()))
			return
		}

		render.Status(r, http.StatusCreated)
		render.Render(w, r, newCreateSessionResponse(cfg, *session))
	}
}

func ConnectSession(cfg *config.Config, p *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := chi.URLParam(r, "session")

		if session == "" {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("missing required path param", "session"))
			return
		}

		var sessionUUID pgtype.UUID
		err := sessionUUID.Scan(session)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("session is not a uuid", err.Error()))
			return
		}

		q := db.New(p)
		dbSession, err := q.GetSessionBySessionId(r.Context(), sessionUUID)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("Error getting session", err.Error()))
		}

		sessionSvc, err := services.NewSessionService(
			r.Context(),
			cfg,
			p,
			services.WithInstance(dbSession.CodechatInstcanceToken, dbSession.CodechatInstance),
		)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("error creating service", err.Error()))
			return
		}
		base64, err := sessionSvc.ConnectSession()
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("error connecting instance", err.Error()))
			return
		}

		render.Status(r, http.StatusOK)
		render.Render(w, r, newConnectSessionResponse(*base64))
	}
}

func StatusSession(cfg *config.Config, p *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := chi.URLParam(r, "session")

		if session == "" {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("missing required path param", "session"))
			return
		}

		var sessionUUID pgtype.UUID
		err := sessionUUID.Scan(session)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("session is not a uuid", err.Error()))
			return
		}

		q := db.New(p)
		dbSession, err := q.GetSessionBySessionId(r.Context(), sessionUUID)
		if err != nil && err.Error() == "no rows in result set" {
			render.Status(r, http.StatusNotFound)
			render.Render(w, r, dto.NewAPIErrorResponse("session not found", ""))
			return
		}
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("Error getting session", err.Error()))
			return
		}

		sessionSvc, err := services.NewSessionService(
			r.Context(),
			cfg,
			p,
			services.WithInstance(dbSession.CodechatInstcanceToken, dbSession.CodechatInstance),
		)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("error creating service", err.Error()))
			return
		}

		resp, err := sessionSvc.FetchInstance()
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("Error fetching instance", err.Error()))
			return
		}
		render.Status(r, http.StatusOK)
		render.Render(
			w,
			r,
			newStatusSessionResponse(cfg, dbSession, resp.ConnectionStatus),
		)
	}
}

func DeleteSession(cfg *config.Config, p *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := chi.URLParam(r, "session")

		if session == "" {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("missing required path param", "session"))
			return
		}

		var sessionUUID pgtype.UUID
		err := sessionUUID.Scan(session)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("session is not a uuid", err.Error()))
			return
		}

		q := db.New(p)
		dbSession, err := q.GetSessionBySessionId(r.Context(), sessionUUID)
		if err != nil && err.Error() == "no rows in result set" {
			render.Status(r, http.StatusNotFound)
			render.Render(w, r, dto.NewAPIErrorResponse("session not found", ""))
			return
		}
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("Error getting session", err.Error()))
			return
		}

		sessionSvc, err := services.NewSessionService(
			r.Context(),
			cfg,
			p,
			services.WithInstance(dbSession.CodechatInstcanceToken, dbSession.CodechatInstance),
		)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("error creating service", err.Error()))
			return
		}

		err = sessionSvc.DeleteInstance()
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("Error deleting instance", err.Error()))
			return
		}

		err = q.DeleteSessionBySessionId(r.Context(), sessionUUID)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.Render(w, r, dto.NewAPIErrorResponse("Error deleting instance", err.Error()))
			return
		}
		render.Status(r, http.StatusNoContent)
	}
}
