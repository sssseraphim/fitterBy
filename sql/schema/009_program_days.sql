-- +goose Up
CREATE TABLE program_days (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
name TEXT NOT NULL,
description TEXT NOT NULL DEFAULT '',
day_order INTEGER NOT NULL,
created_at TIMESTAMP DEFAULT NOW(),
UNIQUE(program_id, day_order));
CREATE INDEX idx_program_days_program ON program_days(program_id);


-- +goose Down
DROP TABLE program_days;
