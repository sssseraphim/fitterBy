-- +goose Up
CREATE TABLE programs(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name TEXT NOT NULL,
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
description TEXT NOT NULL,
media_urls TEXT[],
visibility VARCHAR NOT NULL DEFAULT 'public',
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW());

-- +goose Down
DROP TABLE programs;
