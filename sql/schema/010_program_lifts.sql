-- +goose Up
CREATE TABLE program_lifts (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
program_day_id UUID NOT NULL REFERENCES program_days(id) ON DELETE CASCADE,
exercise_id UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
description TEXT NOT NULL DEFAULT '',
lift_order INTEGER NOT NULL,
sets INTEGER NOT NULL,
reps INTEGER NOT NULL,
created_at TIMESTAMP DEFAULT NOW(),
UNIQUE(program_day_id, lift_order));
CREATE INDEX idx_program_lifts_day ON program_lifts(program_day_id);
CREATE INDEX idx_program_lifts_exercise ON program_lifts(exercise_id);


-- +goose Down
DROP TABLE program_lifts;
