package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/yuqzii/codeforces-insights/internal/codeforces"
)

var ErrNoProblemsForContest = errors.New("there are no stored problems for this contest")
var ErrTooManyProblems = errors.New("amount of problems exceeds alphabet")

func (db *db) GetProblemsWithTags(ctx context.Context, tags []string, minRat, maxRat int) (
	[]codeforces.Problem, error) {

	return db.GetProblemsWithTagsTx(ctx, db.q, tags, minRat, maxRat)
}

func (db *db) GetProblemsWithTagsTx(ctx context.Context, q Querier, tags []string, minRat, maxRat int) (
	[]codeforces.Problem, error) {

	rows, err := q.Query(ctx, `
		SELECT
			p.name,
			p.index,
			p.rating,
			p.tags,
			c.contest_id
		FROM problems p
		JOIN contests c ON p.contest_id = c.id
		WHERE 
			p.tags && $1 AND
			p.rating >= $2 AND p.rating <= $3`,
		tags, minRat, maxRat,
	)
	if err != nil {
		return nil, fmt.Errorf("querying problems: %w", err)
	}

	problems, err := pgx.CollectRows(rows, pgx.RowToStructByName[codeforces.Problem])
	if err != nil {
		return nil, fmt.Errorf("collecting rows: %w", err)
	}

	for i := range problems {
		problems[i].Index = strings.TrimSpace(problems[i].Index)
	}

	return problems, nil
}

// @param id External Codeforces ID of the contest.
func (db *db) GetProblemsFromContest(ctx context.Context, id int) ([]codeforces.Problem, error) {
	return db.GetProblemsFromContestTx(ctx, db.q, id)
}

func (db *db) GetProblemsFromContestTx(ctx context.Context, q Querier, id int) ([]codeforces.Problem, error) {
	rows, err := q.Query(ctx, `
		WITH reference_contest AS (
			SELECT COALESCE(div, 0) AS div, start_time, contest_id
			FROM contests
			WHERE contest_id = $1
		)
		SELECT
			p.name,
			p.index,
			COALESCE(p.rating, 0) AS rating,
			p.tags,
			c.contest_id,
			COALESCE(c.div, 0) AS div
		FROM problems p
		JOIN contests c ON p.contest_id = c.id
		CROSS JOIN reference_contest rc
		WHERE
			c.start_time = rc.start_time AND
			COALESCE(c.div, 0) <= rc.div
		ORDER BY COALESCE(c.div, 0) DESC, p.index ASC`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("querying problems: %w", err)
	}

	problems, err := pgx.CollectRows(rows, pgx.RowToStructByName[probWithDiv])
	if err != nil {
		return nil, fmt.Errorf("collecting rows: %w", err)
	}

	if len(problems) == 0 {
		return nil, ErrNoProblemsForContest
	}

	return correctProblemIndices(problems)
}

func (db *db) UpsertProblem(ctx context.Context, p *codeforces.Problem) error {
	return db.UpsertProblemTx(ctx, db.q, p)
}

func (db *db) UpsertProblemTx(ctx context.Context, q Querier, p *codeforces.Problem) error {
	_, err := q.Exec(ctx, `
		INSERT INTO problems (contest_id, index, name, rating, tags)
		VALUES (
			(SELECT id FROM contests WHERE contest_id = $1),
			$2, $3, NULLIF($4, 0), $5
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

// @return Amount of problems inserted or updated.
func (db *db) UpsertProblemsBatch(ctx context.Context, probs []codeforces.Problem) (int64, error) {
	return db.UpsertProblemsBatchTx(ctx, db.q, probs)
}

func (db *db) UpsertProblemsBatchTx(ctx context.Context, q Querier, probs []codeforces.Problem) (int64, error) {
	batch := &pgx.Batch{}

	query := `
		INSERT INTO problems (contest_id, index, name, rating, tags)
		SELECT c.id, $2, $3, NULLIF($4, 0), $5
		FROM contests c WHERE c.contest_id = $1
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

	var affected int64
	for i := range probs {
		cmd, err := res.Exec()
		if err != nil {
			return 0, fmt.Errorf("batch item %d: %w", i, err)
		}
		affected += cmd.RowsAffected()
	}

	return affected, nil
}

type probWithDiv struct {
	codeforces.Problem
	Div int `db:"div"`
}

// Converts []probWithDiv to []codeforces.Problem and updates their indices,
// in case of multiple contests sharing problems.
// @param probs Slice of probWithDiv, must be sorted by div descending, and then by index ascending.
func correctProblemIndices(probs []probWithDiv) ([]codeforces.Problem, error) {
	res := make([]codeforces.Problem, 0, len(probs))
	var increment uint8 = 0

	for i, p := range probs {
		newProb := p.Problem
		newProb.Index = strings.TrimSpace(newProb.Index)

		if err := updateIndex(&newProb, increment); err != nil {
			return nil, err
		}

		if i < len(probs)-1 && probs[i+1].Div < p.Div {
			// Next problem is from a new contest, update increment.
			increment = newProb.Index[0] - 'A' + 1
			if increment >= 26 {
				return nil, ErrTooManyProblems
			}
		}

		res = append(res, newProb)
	}

	return res, nil
}

func updateIndex(p *codeforces.Problem, increment uint8) error {
	// Early check for large increment to avoid possible overflow when adding later.
	if increment >= 26 {
		return ErrTooManyProblems
	}

	b := []byte(p.Index)
	b[0] += increment
	if b[0] > 'Z' {
		return ErrTooManyProblems
	}

	p.Index = string(b)

	return nil
}
