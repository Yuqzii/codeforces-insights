BEGIN;

DROP INDEX IF EXISTS idx_contests_contest_id;
DROP INDEX IF EXISTS idx_contest_results_contest_id;
DROP INDEX IF EXISTS idx_contest_result_handles_result_id;

COMMIT;
