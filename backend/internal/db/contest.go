package db

import (
	"context"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

func (db *db) ContestExists(ctx context.Context, id int) (exists bool, err error) {
	err = db.conn.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM contests WHERE contest_id=$1)`, id,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, err
}

func (db *db) UpsertContest(ctx context.Context, c *codeforces.Contest) (id int, err error) {
	err = db.conn.QueryRow(ctx, `
		INSERT INTO contests (contest_id, name, start_time, duration) VALUES ($1, $2, $3, $4)
		ON CONFLICT (contest_id) DO UPDATE SET
			name = EXCLUDED.name,
			start_time = EXCLUDED.start_time,
			duration = EXCLUDED.duration
		RETURNING id`,
		c.ID, c.Name, c.StartTime, c.Duration).Scan(&id)
	return id, err
}
