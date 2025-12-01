-- +goose Up
CREATE TABLE refresh_tokens(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
created_at TIMESTAMP DEFAULT NOW(),
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
token_hash TEXT NOT NULL UNIQUE,
expires_at TIMESTAMP NOT NULL);
CREATE INDEX idx_refresh_token_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);

-- +goose Down
DROP TABLE refresh_tokens;
