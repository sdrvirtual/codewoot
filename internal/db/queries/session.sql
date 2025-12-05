
-- name: GetSession :one
SELECT * FROM codechat_session
WHERE id = $1 LIMIT 1;

-- name: ListSessions :many
SELECT * FROM codechat_session;

-- name: GetSessionBySessionId :one
SELECT * FROM codechat_session
WHERE session_id = $1 LIMIT 1;

-- name: CreateSession :one
INSERT INTO codechat_session (
    session_id,
    codechat_instance,
    codechat_instcance_token,
    chatwoot_token,
    chatwoot_account_id,
    chatwoot_inbox_id
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: DeleteSession :exec
DELETE FROM codechat_session
WHERE id = $1;

-- name: DeleteSessionBySessionId :exec
DELETE FROM codechat_session
WHERE session_id = $1;

-- name: UpdateSession :exec
UPDATE codechat_session
  set session_id = $2,
  codechat_instance = $3,
  codechat_instcance_token = $4,
  chatwoot_token = $5,
  chatwoot_account_id = $6,
  chatwoot_inbox_id = $7
WHERE id =  $1;
