-- +goose Up
CREATE TABLE workouts (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL,
program_day_id UUID NOT NULL REFERENCES program_days(id) ON DELETE CASCADE,
created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE workouts;
