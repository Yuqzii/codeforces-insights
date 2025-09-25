package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

func (db *db) InsertContestResults(ctx context.Context, contestants []codeforces.Contestant, id int) error {
	return db.InsertContestResultsTx(ctx, db.q, contestants, id)
}

func (db *db) InsertContestResultsTx(ctx context.Context, q Querier,
	contestants []codeforces.Contestant, id int) error {

	batch := &pgx.Batch{}
	for _, c := range contestants {
		batch.Queue(`
			WITH new_result AS (
				INSERT INTO contest_results (contest_id, rank, old_rating, new_rating, points)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id
			)
			INSERT INTO contest_result_handles (contest_result_id, handle)
			SELECT new_result.id, UNNEST($6::varchar(32)[])
			FROM new_result`,
			id, c.Rank, c.OldRating, c.NewRating, c.Points, c.MemberHandles)
	}
	br := q.SendBatch(ctx, batch)
	if err := br.Close(); err != nil {
		return fmt.Errorf("closing batch result: %w", err)
	}

	return nil
}
