-- name: CreateProgram :one
INSERT INTO programs (name, user_id, description, media_urls, visibility)
VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
		)
RETURNING *;

-- name: CreateProgramDay :one
INSERT INTO program_days (program_id, name, description, day_order)
VALUES (
		$1,
		$2,
		$3,
		$4
		)
RETURNING *;

-- name: CreateProgramLift :one
INSERT INTO program_lifts (program_day_id, exercise_id, description, lift_order, sets, reps)
VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6
		)
RETURNING *;


-- name: GetProgram :one
SELECT programs.*, users.name as author_name
FROM programs
LEFT JOIN users ON programs.user_id = users.id
WHERE programs.id = $1;

-- name: GetPrograms :many
SELECT programs.*, users.name as author_name
FROM programs 
LEFT JOIN users 
ON programs.user_id = users.id
ORDER BY programs.created_at DESC
LIMIT 50;

-- name: GetProgramDays :many
SELECT *
FROM program_days
WHERE program_id = $1
ORDER BY day_order ASC;

-- name: GetProgramDayLifts :many
SELECT p.*
FROM program_lifts p
LEFT JOIN exercises e ON p.exercise_id = e.id
WHERE program_day_id = $1
ORDER BY lift_order ASC;

-- name: GetUserSubscribedPrograms :many
SELECT * FROM users_programs
WHERE user_id = $1;

-- name: SubscribeToProgram :exec
INSERT INTO users_programs(
		user_id,
		program_id
) VALUES (
		$1,
		$2
		);
