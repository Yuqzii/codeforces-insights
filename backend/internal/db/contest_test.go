package db

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

func setupMockDB(t *testing.T) (pgxmock.PgxPoolIface, *db) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to set up mock DB: %v", err)
	}

	return mock, &db{q: mock}
}

func TestContestExists(t *testing.T) {
	ctx := context.Background()
	mock, db := setupMockDB(t)

	t.Run("Contest exists", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"exists"})
		rows.AddRow(true)
		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM contests WHERE contest_id=\$1\)`).
			WithArgs(69).
			WillReturnRows(rows)

		exists, err := db.ContestExists(ctx, 69)
		assert.Nil(t, err, err)
		assert.True(t, exists, "expected contest to exist")
	})

	t.Run("Contest does not exist", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM contests WHERE contest_id=\$1\)`).
			WithArgs(99).
			WillReturnRows(rows)

		exists, err := db.ContestExists(ctx, 99)
		assert.Nil(t, err, err)
		assert.False(t, exists, "expected contest to not exist")
	})

	t.Run("Query error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM contests WHERE contest_id=\$1\)`).
			WithArgs(13).
			WillReturnError(errors.New("query failed"))

		exists, err := db.ContestExists(ctx, 13)
		assert.NotNil(t, err, "expected error")
		assert.False(t, exists, "expected exists=false on error")
	})
}

func TestContestsExists(t *testing.T) {
	ctx := context.Background()
	mock, db := setupMockDB(t)

	t.Run("All contests exist", func(t *testing.T) {
		input := []int{1, 2, 50}
		rows := pgxmock.NewRows([]string{"contest_id"}).
			AddRow(1).AddRow(2).AddRow(50)
		mock.ExpectQuery(`SELECT contest_id FROM contests WHERE contest_id = ANY\(\$1\)`).
			WithArgs(input).
			WillReturnRows(rows)

		existing, err := db.ContestsExists(ctx, input)
		assert.Nil(t, err, err)
		for _, val := range input {
			_, ok := existing[val]
			assert.Truef(t, ok, "expected contest %d to exist", val)
		}
	})

	t.Run("Some contests exist", func(t *testing.T) {
		input := []int{42, 1, 69, 100}
		shouldExist := []bool{true, false, true, false}
		rows := pgxmock.NewRows([]string{"contest_id"}).
			AddRow(42).AddRow(69)
		mock.ExpectQuery(`SELECT contest_id FROM contests WHERE contest_id = ANY\(\$1\)`).
			WithArgs(input).
			WillReturnRows(rows)

		existing, err := db.ContestsExists(ctx, input)
		assert.Nil(t, err, err)
		for i := range input {
			_, ok := existing[input[i]]
			assert.Equal(t, shouldExist[i], ok)
		}
	})

	t.Run("No contests exist", func(t *testing.T) {
		input := []int{10, 1337, 3, 69420}
		rows := pgxmock.NewRows([]string{"contest_id"})
		mock.ExpectQuery(`SELECT contest_id FROM contests WHERE contest_id = ANY\(\$1\)`).
			WithArgs(input).
			WillReturnRows(rows)

		existing, err := db.ContestsExists(ctx, input)
		assert.Nil(t, err, err)
		for _, val := range input {
			_, ok := existing[val]
			assert.Falsef(t, ok, "did not expect contest %d to exist", val)
		}
	})

	t.Run("Query error", func(t *testing.T) {
		input := []int{10, 20}
		mock.ExpectQuery(`SELECT contest_id FROM contests WHERE contest_id = ANY\(\$1\)`).
			WithArgs(input).
			WillReturnError(errors.New("query failed"))

		existing, err := db.ContestsExists(ctx, input)
		assert.NotNil(t, err, "expected error")
		assert.Nilf(t, existing, "expected existing=nil on error, got %v", existing)
	})
}

func TestUpsertContestsTx(t *testing.T) {
	ctx := context.Background()
	mock, db := setupMockDB(t)

	c := &codeforces.Contest{
		ID:        1337,
		Name:      "Mock Contest",
		StartTime: time.Date(2025, 3, 10, 0, 0, 0, 0, time.UTC),
		Duration:  7200,
	}

	t.Run("Successful upsert", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id"}).AddRow(42)
		mock.ExpectQuery(`INSERT INTO contests`).
			WithArgs(c.ID, c.Name, c.StartTime, c.Duration).
			WillReturnRows(rows)

		id, err := db.UpsertContestTx(ctx, mock, c)
		assert.Nil(t, err, err)
		assert.Equalf(t, id, 42, "expected id 42, got %d", id)
	})

	t.Run("Query error", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO contests`).
			WithArgs(c.ID, c.Name, c.StartTime, c.Duration).
			WillReturnError(errors.New("insert failed"))

		_, err := db.UpsertContestTx(ctx, mock, c)
		assert.NotNil(t, err)
	})
}
