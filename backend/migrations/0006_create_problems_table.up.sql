BEGIN;

CREATE TABLE IF NOT EXISTS problems (
	id SERIAL PRIMARY KEY,
	contest_id INT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
	index CHAR(2) NOT NULL,
	name VARCHAR(256) NOT NULL,
	rating INT,
	tags TEXT[] NOT NULL DEFAULT '{}',
	updated_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE problems ADD CONSTRAINT unique_problem_per_contest UNIQUE (contest_id, index);

CREATE TRIGGER set_problems_updated_at
BEFORE UPDATE ON problems
FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();

CREATE INDEX IF NOT EXISTS idx_problems_tags ON problems USING GIN (tags);

CREATE INDEX IF NOT EXISTS idx_problems_contest_id ON problems(contest_id);

COMMIT;
