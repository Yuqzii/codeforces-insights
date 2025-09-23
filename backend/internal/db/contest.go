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

func (db *db) ContestsExists(ctx context.Context, ids []int) (existingIDs map[int]struct{}, err error) {
	rows, err := db.conn.Query(ctx, `
		SELECT contest_id FROM contests WHERE contest_id = ANY($1)`,
		ids,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existingIDs = make(map[int]struct{})
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		existingIDs[id] = struct{}{}
	}

	return existingIDs, nil
}

func (db *db) UpsertContest(ctx context.Context, c *codeforces.Contest) (id int, err error) {
	return db.UpsertContestTx(ctx, db.conn, c)
}

func (db *db) UpsertContestTx(ctx context.Context, q Querier, c *codeforces.Contest) (id int, err error) {
	err = q.QueryRow(ctx, `
		INSERT INTO contests (contest_id, name, start_time, duration) VALUES ($1, $2, $3, $4)
		ON CONFLICT (contest_id) DO UPDATE SET
			name = EXCLUDED.name,
			start_time = EXCLUDED.start_time,
			duration = EXCLUDED.duration
		RETURNING id`,
		c.ID, c.Name, c.StartTime, c.Duration,
	).Scan(&id)
	return id, err
}
