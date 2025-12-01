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

-- name: FollowUser :exec
INSERT INTO user_follows(follower_id, followed_id)
VALUES (
		$1,
		$2
		);

-- name: GetFollowedUsers :many
SELECT users.id, users.name, users.created_at, users.updated_at, users.bio, users.premium
FROM users
INNER JOIN user_follows ON users.id = user_follows.followed_id
WHERE user_follows.follower_id = $1;
