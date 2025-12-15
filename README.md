# Codechat ↔ Chatwoot Relay

A Go service that bridges Codechat and Chatwoot. It ingests webhooks, manages sessions, validates phone numbers, and transcodes audio to ensure compatibility across platforms.

## Features

- Webhook handlers for Chatwoot and Codechat with strict JSON Content-Type and 2MB body limits
- Session management persisted in PostgreSQL via pgx, with typed models and queries
- Audio transcoding from OGG to MP3 using ffmpeg (libmp3lame)
- Phone number validation for Brazilian and international formats
- Strongly-typed DTOs for both external APIs
- Configuration via environment variables and `.env` files (godotenv)

## Architecture Overview

- `internal/handlers/*`: HTTP handlers
  - `ChatwootWebhook`, `CodechatWebhook`, `CreateSession`
- `internal/services/*`: business logic
  - `RelayService`, `SessionService`
- `internal/codechat/*` and `internal/chatwoot/*`: API clients
  - Shared `Option` pattern and `newRequest` helpers supporting `io.Reader` bodies or JSON
- `internal/audio/transcoder.go`: OGG → MP3 transcoder (ffmpeg piping: `pipe:0` → `pipe:1` with `libmp3lame`)
- `internal/db/*`: models and queries (pgx/pgxpool, generated `session.sql.go`)
- `internal/dto/*`: typed payloads for Chatwoot and Codechat webhooks
- `internal/utils/phone.go`: phone normalization and validation
- `internal/config/config.go`: environment-driven configuration (via `godotenv`)
- `internal/server/server.go`: HTTP server helpers (e.g., body capture writer)

## Prerequisites

- Go 1.20+ (recommended)
- ffmpeg installed with `libmp3lame` codec
- PostgreSQL (pgx/pgxpool used by the app)
- Chatwoot and Codechat credentials as required for your deployment

## Configuration

Environment variables read by `internal/config/config.go`:
- `PORT`: server port (default: `8080`)
- `HOST`: server host (default: `0.0.0.0`)
- `API_URL`: public base URL of this service (default: `http://localhost:8080`)
- `CHATWOOT_URL`: base URL of your Chatwoot deployment

Additional environment variables used by the project:
- `CODECHAT_URL`: base URL of your Codechat deployment
- `CODECHAT_KEY`: API key/token for Codechat
- `DB_URL`: PostgreSQL connection string (e.g., `postgresql://user:pass@host/db`)
- `GOOSE_DRIVER`: database driver for migrations (e.g., `postgres`)
- `GOOSE_DBSTRING`: connection string used by goose (usually same as `DB_URL`)
- `GOOSE_MIGRATION_DIR`: directory of migration files (e.g., `./internal/db/migrations`)

The project loads `.env` automatically via `godotenv`. Create a `.env` file in the repository root and set the variables above as needed. You can copy `.env.example` to `.env` and adjust values for your environment.

Note: Database connection variables are expected to be wired where the `pgxpool.Pool` is created (outside of `internal/config/config.go`). Configure your Postgres connection accordingly in your server/bootstrap code.

## Running

Wire the HTTP routes in your main package to the handlers:
- `internal/handlers.ChatwootWebhook(cfg, pool)`
- `internal/handlers.CodechatWebhook(cfg, pool)`
- `internal/handlers.CreateSession(cfg, pool)`

Ensure your server constructs:
- A loaded `config.Config` (via `config.Load()`)
- A `pgxpool.Pool` connected to your PostgreSQL instance

Start the server using your main package entry point (commonly under `cmd/`), then send requests to the configured webhook and session endpoints.

## API Overview

### Create Session
JSON payload shape (from `internal/dto/session.go`):
```json
{
  "session_id": "b1f2f0f4-7930-4e90-9b6f-8c0f3c9b6c12",
  "description": "Support instance for WhatsApp",
  "chatwoot": {
    "inbox_id": 123,
    "account_id": 456,
    "token": "chatwoot_api_token_here"
  }
}
```
- If `session_id` is omitted, the service generates a UUID.
- The service persists the session and returns identifiers/tokens as implemented in `SessionService`.

### Webhooks
- Chatwoot: conforms to `internal/dto/chatwoot.go` (`ChatwootWebhook`). Handler enforces `Content-Type: application/json`.
- Codechat: conforms to `internal/dto/codechat.go` (`CodechatWebhook`). Handler enforces `Content-Type: application/json`.

## Quickstart

1. Ensure dependencies installed:
   - Go 1.20+ (or newer)
   - ffmpeg with `libmp3lame` codec available
   - PostgreSQL accessible from your development environment

2. Configure environment:
   - Copy `.env.example` to `.env` and adjust values for your environment.
   - Ensure `PORT`, `HOST`, `API_URL`, `CHATWOOT_URL`, `CODECHAT_URL`, `CODECHAT_KEY`, and `DB_URL` are set.
   - Provide your PostgreSQL connection configuration where you construct the `pgxpool.Pool`.

3. Run the service:
   - Load configuration via `config.Load()`.
   - Initialize a `pgxpool.Pool`.
   - Wire HTTP routes to the handlers:
     - `internal/handlers.ChatwootWebhook(cfg, pool)`
     - `internal/handlers.CodechatWebhook(cfg, pool)`
     - `internal/handlers.CreateSession(cfg, pool)`
   - Start your HTTP server with the configured routes.

## Roadmap

- Chatwoot
  - [x] Support text messages
  - [x] Support audio messages
  - [x] Support document messages

- Codechat
  - [x] Support text messages
  - [x] Support audio messages
  - [ ] Support media messages
  - [ ] Support document messages
  - [ ] Support pools
  - [ ] Support contact messages
  - [ ] Support location messages
  - [ ] Support live location messages
  - [ ] Support pvt messages
  - [ ] Support interactive messages
  - [ ] Support sticker messages


- Platform
  - [ ] Handle message read events

## License

TBD. Add a LICENSE file appropriate for your distribution.
