-- +goose Up
CREATE TABLE user_follows(
follower_id UUID REFERENCES users(id),
followed_id UUID REFERENCES users(id),
created_at TIMESTAMP DEFAULT NOW(),
PRIMARY KEY (follower_id, followed_id)
);

-- +goose Down
DROP TABLE user_follows;
