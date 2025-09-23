package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

func (db *db) InsertContestResults(ctx context.Context, contestants []codeforces.Contestant, id int) error {
	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	for _, c := range contestants {
		_, err := tx.Exec(ctx, `
			INSERT INTO contest_results
			(contest_id, rank, old_rating, new_rating, points) VALUES ($1, $2, $3, $4, $5)`,
			id, c.Rank, c.OldRating, c.NewRating, c.Points)
		if err != nil {
			return fmt.Errorf("inserting contest result: %w", err)
		}
	}

	return nil
}
