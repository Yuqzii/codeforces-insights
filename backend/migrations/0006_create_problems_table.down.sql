BEGIN;

DROP INDEX IF EXISTS idx_problems_contest_id;
DROP INDEX IF EXISTS idx_problems_tags;

DROP TRIGGER IF EXISTS set_problems_updated_at ON problems;

DROP TABLE IF EXISTS problems;

COMMIT;
