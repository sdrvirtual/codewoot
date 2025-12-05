-- +goose Up
-- +goose StatementBegin
CREATE TABLE codechat_session (
       id SERIAL PRIMARY KEY,
       session_id UUID NOT NULL UNIQUE,
       codechat_instance VARCHAR(255) NOT NULL,
       codechat_instcance_token VARCHAR(1024) NOT NULL,
       chatwoot_token VARCHAR(1024) NOT NULL,
       chatwoot_account_id int NOT NULL,
       chatwoot_inbox_id int NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER update_products_updated_at
BEFORE UPDATE ON codechat_session
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE codechat_session;
DROP FUNCTION update_modified_column;
-- +goose StatementEnd
