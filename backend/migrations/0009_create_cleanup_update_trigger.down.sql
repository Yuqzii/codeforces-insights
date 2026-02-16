BEGIN;

DROP FUNCTION IF EXISTS delete_contest_dependents;
DROP TRIGGER IF EXISTS delete_contest_dependents_on_update;

COMMIT;
