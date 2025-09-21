BEGIN;

CREATE TABLE IF NOT EXISTS contest_results (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	contest_id INT NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
	rank INT NOT NULL,
	old_rating INT NOT NULL,
	new_rating INT NOT NULL,
	points REAL NOT NULL,
	updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TRIGGER set_contest_results_updated_at
BEFORE UPDATE ON contest_results
FOR EACH ROW
EXECUTE PROCEDURE set_updated_at();

COMMIT;
