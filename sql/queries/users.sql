-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name, email, bio, hashed_password)
VALUES (
		gen_random_uuid(),
		NOW(),
		NOW(),
		$1,
		$2,
		$3,
		$4
)
RETURNING *;

-- name: GetUser :one
SELECT id, created_at, updated_at, name, email, bio, premium
FROM users
WHERE id = $1;

-- name: GetUserLogin :one
SELECT id, created_at, updated_at, name, email, bio, premium, hashed_password
FROM users
WHERE email = $1;

-- name: UpdateBio :exec
UPDATE users
SET bio = $2,
updated_at = Now()
WHERE id = $1;
