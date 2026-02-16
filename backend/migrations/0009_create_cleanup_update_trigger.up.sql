BEGIN;

CREATE OR REPLACE FUNCTION delete_contest_dependents()
RETURNS TRIGGER AS $$
BEGIN
	-- When the fetcher updates a contest it also inserts the contest results again.
	-- This prevents keeping invalid contest results.
	DELETE FROM contest_results WHERE contest_id = OLD.id;
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER delete_contest_dependents_on_update
BEFORE UPDATE ON contests
FOR EACH ROW
EXECUTE FUNCTION delete_contest_dependents();

COMMIT;
