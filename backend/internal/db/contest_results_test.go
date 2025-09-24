package db

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

func TestInsertContestResultsTx(t *testing.T) {
	ctx := context.Background()
	mock, db := setupMockDB(t)

	contestants := []codeforces.Contestant{
		{
			Rank:          1,
			OldRating:     1337,
			NewRating:     1600,
			Points:        3400,
			MemberHandles: []string{"julian"},
		},
		{
			Rank:          2,
			OldRating:     1400,
			NewRating:     1450,
			Points:        3350,
			MemberHandles: []string{"martin, marius"},
		},
		{
			Rank:          3,
			OldRating:     1248,
			NewRating:     1150,
			Points:        2500,
			MemberHandles: []string{"gru"},
		},
	}

	t.Run("Successful batch", func(t *testing.T) {
		eb := mock.ExpectBatch()
		id := 67

		for _, c := range contestants {
			eb.ExpectExec(`WITH new_result AS .* INSERT INTO contest_result_handles`).
				WithArgs(id, c.Rank, c.OldRating, c.NewRating, c.Points, c.MemberHandles).
				WillReturnResult(pgxmock.NewResult("INSERT", 1))
		}

		err := db.InsertContestResultsTx(ctx, mock, contestants, id)
		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Batch send fails", func(t *testing.T) {
		eb := mock.ExpectBatch()
		eb.WillReturnError(errors.New("batch send error"))

		id := 700
		for _, c := range contestants {
			eb.ExpectExec(`.*`).
				WithArgs(id, c.Rank, c.OldRating, c.NewRating, c.Points, c.MemberHandles).
				WillReturnResult(pgxmock.NewResult("INSERT", 1))

		}

		err := db.InsertContestResultsTx(ctx, mock, contestants, id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "batch send error")
	})
}
