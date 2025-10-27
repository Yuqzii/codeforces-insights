package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

var ErrContestNotStored = errors.New("did not find contest in db")

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

func (db *db) GetContestResults(ctx context.Context, id int) (
	[]codeforces.Contestant, *codeforces.Contest, error) {

	return db.GetContestResultsTx(ctx, db.q, id)
}

func (db *db) GetContestResultsTx(ctx context.Context, q Querier, id int) (
	[]codeforces.Contestant, *codeforces.Contest, error) {

	// Get contest
	var contest codeforces.Contest
	var internalID int
	err := q.QueryRow(ctx, `
		SELECT name, start_time, duration, contest_id, id FROM contests WHERE contest_id=$1`,
		id,
	).Scan(&contest.Name, &contest.StartTime, &contest.Duration, &contest.ID, &internalID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrContestNotStored
		}
		return nil, nil, fmt.Errorf("querying contests: %w", err)
	}

	// Get contestants with handles
	rows, err := q.Query(ctx, `
		SELECT
			cr.rank,
			cr.old_rating,
			cr.new_rating,
			cr.points,
			cr.id,
			COALESCE(ARRAY_AGG(crh.handle) FILTER (WHERE crh.handle IS NOT NULL), '{}') AS handles
		FROM contest_results AS cr
		LEFT JOIN contest_result_handles crh ON crh.contest_result_id=cr.id
		WHERE cr.contest_id = $1
		GROUP BY cr.rank, cr.old_rating, cr.new_rating, cr.points, cr.id`,
		internalID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("querying contest_results: %w", err)
	}

	contestants, err := scanToContestants(rows)
	if err != nil {
		return nil, nil, fmt.Errorf("scanning into contestant: %w", err)
	}

	return contestants, &contest, nil
}

func scanToContestants(rows pgx.Rows) ([]codeforces.Contestant, error) {
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (codeforces.Contestant, error) {
		var c codeforces.Contestant

		err := row.Scan(
			&c.Rank,
			&c.OldRating,
			&c.NewRating,
			&c.Points,
			&c.ID,
			&c.MemberHandles,
		)

		return c, err
	})
}
