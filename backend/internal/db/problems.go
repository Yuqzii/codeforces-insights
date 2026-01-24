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
			p.tags,
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

func (db *db) UpsertProblem(ctx context.Context, p *codeforces.Problem) error {
	return db.UpsertProblemTx(ctx, db.q, p)
}

func (db *db) UpsertProblemTx(ctx context.Context, q Querier, p *codeforces.Problem) error {
	_, err := q.Exec(ctx, `
		INSERT INTO problems (contest_id, index, name, rating, tags)
		VALUES (
			(SELECT id FROM contests WHERE contest_id = $1),
			$2, $3, $4, $5
		)
		ON CONFLICT (contest_id, index)
		DO UPDATE SET
			name = EXCLUDED.name,
			rating = EXCLUDED.rating,
			tags = EXCLUDED.tags`,
		p.ContestID, p.Index, p.Name, p.Rating, p.Tags,
	)

	return err
}

func (db *db) UpsertProblemsBatch(ctx context.Context, probs []codeforces.Problem) error {
	return db.UpsertProblemsBatchTx(ctx, db.q, probs)
}

func (db *db) UpsertProblemsBatchTx(ctx context.Context, q Querier, probs []codeforces.Problem) error {
	batch := &pgx.Batch{}

	query := `
		INSERT INTO problems (contest_id, index, name, rating, tags)
		VALUES (
			(SELECT id FROM contests WHERE contest_id = $1),
			$2, $3, $4, $5
		)
		ON CONFLICT (contest_id, index)
		DO UPDATE SET
			name = EXCLUDED.name,
			rating = EXCLUDED.rating,
			tags = EXCLUDED.tags`

	for _, p := range probs {
		batch.Queue(query, p.ContestID, p.Index, p.Name, p.Rating, p.Tags)
	}

	res := q.SendBatch(ctx, batch)
	defer res.Close() //nolint:errcheck

	for i := range probs {
		_, err := res.Exec()
		if err != nil {
			return fmt.Errorf("batch item %d: %w", i, err)
		}
	}

	return nil
}
