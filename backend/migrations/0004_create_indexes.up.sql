BEGIN;

CREATE INDEX IF NOT EXISTS idx_contests_contest_id
ON contests(contest_id);

CREATE INDEX IF NOT EXISTS idx_contest_results_contest_id
ON contest_results(contest_id);

CREATE INDEX IF NOT EXISTS idx_contest_result_handles_result_id
ON contest_result_handles(contest_result_id);

COMMIT;
