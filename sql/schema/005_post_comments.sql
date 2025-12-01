-- +goose Up
CREATE TABLE posts_comments(
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID REFERENCES users(id),
post_id UUID REFERENCES posts(id),
content TEXT NOT NULL,
created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE posts_comments;
