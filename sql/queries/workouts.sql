-- name: CreateWorkout :one
INSERT INTO workouts(user_id, program_day_id)
VALUES (
		$1,
		$2
		)
RETURNING *;

-- name: CreateWorkoutLift :one
INSERT INTO users_lifts(user_id, workout_id, exercise_id, weight, sets, reps, lift_order)
VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7)
RETURNING *;

-- name: GetUsersLiftsByExercise :many
SELECT * FROM users_lifts
WHERE exercise_id = $1
ORDER BY created_at DESC;

-- name: GetUsersWorkouts :many
SELECT * FROM workouts
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetWorkoutByID :one
SELECT * FROM workouts
WHERE id = $1;

-- name: GetWorkoutLifts :many
SELECT * FROM users_lifts
WHERE workout_id = $1
ORDER BY lift_order ASC;
