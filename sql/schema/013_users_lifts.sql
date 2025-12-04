-- +goose Up
CREATE TABLE users_lifts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
		workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT NOW(),
		weight INTEGER NOT NULL,
		sets INTEGER NOT NULL,
		reps INTEGER NOT NULL,
		lift_order INTEGER NOT NULL
);
CREATE INDEX idx_users_lifts ON users_lifts(exercise_id);

-- +goose Down
DROP TABLE users_lifts;
