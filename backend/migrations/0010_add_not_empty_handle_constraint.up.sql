ALTER TABLE contest_result_handles
ADD CONSTRAINT not_empty_handle
CHECK (handle <> '');
