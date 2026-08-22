package db

import (
	"context"
	"errors"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yuqzii/codeforces-insights/internal/codeforces"
)

func TestGetContestResultsFromHandle(t *testing.T) {
	ctx := context.Background()
	columns := []string{"rank", "old_rating", "new_rating", "points", "id", "contest_id", "handles"}

	tests := []struct {
		name string
		want []codeforces.Contestant
	}{
		{
			name: "one result",
			want: []codeforces.Contestant{
				{
					Rank:          1,
					OldRating:     2100,
					NewRating:     2200,
					Points:        3400,
					ID:            10,
					ContestID:     100,
					MemberHandles: []string{"tourist"},
				},
			},
		},
		{
			name: "multiple results",
			want: []codeforces.Contestant{
				{
					Rank:          1,
					OldRating:     2100,
					NewRating:     2200,
					Points:        3400,
					ID:            10,
					ContestID:     100,
					MemberHandles: []string{"tourist"},
				},
				{
					Rank:          2,
					OldRating:     2200,
					NewRating:     2250,
					Points:        3200,
					ID:            20,
					ContestID:     200,
					MemberHandles: []string{"tourist", "teammate"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, db := setupMockDB(t)
			rows := pgxmock.NewRows(columns)
			for _, contestant := range tt.want {
				rows.AddRow(
					contestant.Rank,
					contestant.OldRating,
					contestant.NewRating,
					contestant.Points,
					contestant.ID,
					contestant.ContestID,
					contestant.MemberHandles,
				)
			}
			mock.ExpectQuery(`WITH matching_results AS`).
				WithArgs("tourist").
				WillReturnRows(rows)

			got, err := db.GetContestResultsFromHandle(ctx, "Tourist")
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

	t.Run("no results", func(t *testing.T) {
		mock, db := setupMockDB(t)
		mock.ExpectQuery(`WITH matching_results AS`).
			WithArgs("nobody").
			WillReturnRows(pgxmock.NewRows(columns))

		got, err := db.GetContestResultsFromHandle(ctx, "Nobody")
		assert.ErrorIs(t, err, ErrNoResultsWithHandle)
		assert.Nil(t, got)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

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
