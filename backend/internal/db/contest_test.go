package db

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
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
		assert.Nil(t, err)
		assert.True(t, exists, "expected contest to exist")
	})

	t.Run("Contest does not exist", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM contests WHERE contest_id=\$1\)`).
			WithArgs(99).
			WillReturnRows(rows)

		exists, err := db.ContestExists(ctx, 99)
		assert.Nil(t, err)
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
