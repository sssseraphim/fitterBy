-- +goose Up
CREATE TABLE posts_likes(
user_id UUID REFERENCES users(id),
post_id UUID REFERENCES posts(id),
created_at TIMESTAMP DEFAULT NOW(),
PRIMARY KEY (user_id, post_id)
);

-- +goose Down
DROP TABLE posts_likes;
