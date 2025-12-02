-- name: GetExerciseById :one
SELECT exercises.*, users.name as author_name
FROM exercises
LEFT JOIN users ON exercises.user_id = users.id
WHERE exercises.id = $1;

-- name: GetExercises :many
SELECT exercises.*, users.name as author_name
FROM exercises
LEFT JOIN users ON exercises.user_id = users.id;

-- name: CreateExercise :one
INSERT INTO exercises(name, user_id, description, media_urls)
VALUES (
		$1,
		$2,
		$3,
		$4
		)
RETURNING *;

