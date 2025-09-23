BEGIN;

CREATE TABLE IF NOT EXISTS contest_result_handles (
	contest_result_id INT NOT NULL REFERENCES contest_results(id) ON DELETE CASCADE,
	handle VARCHAR(32) NOT NULL
);

COMMIT;
