-- +goose Up
CREATE TABLE users_programs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		program_id UUID NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT NOW(),
		current_day_order INTEGER DEFAULT 1,
		status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'paused', 'completed', 'abandoned')),
		UNIQUE(user_id, program_id)
);
CREATE INDEX idx_users_programs_user ON users_programs(user_id, status);
CREATE INDEX idx_users_programs_program ON users_programs(program_id);

-- +goose Down
DROP TABLE users_programs;
