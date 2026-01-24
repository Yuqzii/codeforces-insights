package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

func (db *db) GetProblemsWithTags(ctx context.Context, tags []string) ([]codeforces.Problem, error) {
	return db.GetProblemsWithTagsTx(ctx, db.q, tags)
}

func (db *db) GetProblemsWithTagsTx(ctx context.Context, q Querier, tags []string) ([]codeforces.Problem, error) {
	rows, err := q.Query(ctx, `
		SELECT
			p.name,
			p.index,
			p.rating,
			p.tags,
			c.contest_id
		FROM problems p
		JOIN contests c ON p.contest_id = c.id
		WHERE p.tags && $1`,
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("querying problems: %w", err)
	}

	problems, err := pgx.CollectRows(rows, pgx.RowToStructByName[codeforces.Problem])
	if err != nil {
		return nil, fmt.Errorf("collecting rows: %w", err)
	}

	return problems, nil
}

// @param id External Codeforces ID of the contest.
func (db *db) GetProblemsFromContest(ctx context.Context, id int) ([]codeforces.Problem, error) {
	return db.GetProblemsFromContestTx(ctx, db.q, id)
}

func (db *db) GetProblemsFromContestTx(ctx context.Context, q Querier, id int) ([]codeforces.Problem, error) {
	rows, err := q.Query(ctx, `
		SELECT
			p.name,
			p.index,
			p.rating,
			p.tags
			c.contest_id
		FROM problems p
		JOIN contests c ON p.contest_id = c.id
		WHERE c.contest_id = $1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("querying problems: %w", err)
	}

	problems, err := pgx.CollectRows(rows, pgx.RowToStructByName[codeforces.Problem])
	if err != nil {
		return nil, fmt.Errorf("collecting rows: %w", err)
	}

	return problems, nil
}
