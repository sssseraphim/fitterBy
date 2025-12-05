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
SELECT posts.id, users.id as author_id, users.name as author_name, posts.created_at, posts.visibility, posts.media_urls, posts.content, 
		(SELECT COUNT(*) FROM posts_likes where posts.id = posts_likes.post_id) as like_count,
		(SELECT COUNT(*) FROM posts_comments where posts.id = posts_comments.post_id) as comments_count
FROM posts
LEFT JOIN users ON posts.user_id = users.id
WHERE posts.id = $1;

-- name: GetFollowedPosts :many
SELECT posts.id, users.id as author_id, users.name as author_name, posts.created_at, posts.visibility, posts.media_urls, posts.content,
		(SELECT COUNT(*) FROM posts_likes where posts.id = posts_likes.post_id) as like_count,
		(SELECT COUNT(*) FROM posts_comments where posts.id = posts_comments.post_id) as comments_count,
		EXISTS(SELECT 1 FROM posts_likes WHERE posts_likes.post_id = posts.id AND posts_likes.user_id = $1) as user_liked
FROM posts 
INNER JOIN user_follows ON posts.user_id = user_follows.followed_id
INNER JOIN users ON posts.user_id = users.id
WHERE user_follows.follower_id = $1
AND posts.visibility IN ('public', 'followers')
ORDER BY posts.created_at DESC;

-- name: GetLikeCount :one
SELECT COUNT(*)
FROM posts_likes
WHERE post_id = $1;

-- name: LikePost :exec
INSERT INTO posts_likes(user_id, post_id)
VALUES (
		$1,
		$2
		);

-- name: CommentOnPost :exec
INSERT INTO posts_comments(user_id, post_id, content)
VALUES (
		$1,
		$2,
		$3
		);

-- name: GetPostComments :many
SELECT posts_comments.*, users.name as commenter_name
FROM posts_comments
LEFT JOIN users ON posts_comments.user_id = users.id
WHERE posts_comments.post_id = $1
ORDER BY posts_comments.created_at
LIMIT 50;

-- name: CheckUserLikedPost :one
SELECT EXISTS (
		SELECT 1
		FROM posts_likes
		WHERE post_id = $1
		AND user_id = $2
) AS user_liked;

-- name: DeleteUserPostLike :exec
DELETE FROM posts_likes
WHERE post_id = $1
AND user_id = $2;

