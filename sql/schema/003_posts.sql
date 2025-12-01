-- +goose Up
CREATE TABLE posts(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
content TEXT NOT NULL,
media_urls TEXT[],
visibility VARCHAR NOT NULL DEFAULT 'public',
like_count INTEGER NOT NULL DEFAULT 0,
comment_count INTEGER NOT NULL DEFAULT 0,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW());

CREATE INDEX idx_poster_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);

-- +goose Down
DROP TABLE posts;
