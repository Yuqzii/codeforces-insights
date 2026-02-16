BEGIN;

DROP TRIGGER IF EXISTS delete_contest_dependents_on_update ON contests;
DROP FUNCTION IF EXISTS delete_contest_dependents;

COMMIT;
