package db

import "context"

func (db *db) ContestExists(ctx context.Context, id int) (exists bool, err error) {
	err = db.conn.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM contests WHERE contest_id=$1)`, id,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}
