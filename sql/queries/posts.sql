-- name: CreatePost :one
INSERT INTO posts(user_id, content, media_urls, visibility)
VALUES (
		$1,
		$2,
		$3,
		$4
		)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1;

-- name: GetFollowedPosts :many
SELECT posts.id, users.id as author_id, users.name as author_name, posts.created_at, posts.visibility, posts.media_urls, posts.content
FROM posts
INNER JOIN user_follows ON posts.user_id = user_follows.followed_id
INNER JOIN users ON post.user_id = users.id
WHERE user_follows.follower_id = $1
AND posts.visibility IN ('public', 'followers')
ORDER BY posts.created_at DESC;
