-- +goose Up
CREATE TABLE users(
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
name TEXT NOT NULL,
email TEXT NOT NULL UNIQUE,
bio TEXT NOT NULL,
hashed_password TEXT NOT NULL,
premium BOOLEAN DEFAULT false);


-- +goose Down
DROP TABLE users;
