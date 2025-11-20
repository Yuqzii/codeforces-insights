package store

import (
	"context"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type API interface {
	GetContestStandings(ctx context.Context, id int) (
		[]codeforces.Contestant, *codeforces.Contest, error)
	GetContestRatingChanges(ctx context.Context, id int) ([]codeforces.RatingChange, error)
	GetRatingChanges(ctx context.Context, handle string) ([]codeforces.RatingChange, error)
}

type DB interface {
	GetContestResults(ctx context.Context, id int, idIsInternal bool) (
		[]codeforces.Contestant, *codeforces.Contest, error)
	GetContestResultsFromHandle(ctx context.Context, handle string) (
		[]codeforces.Contestant, error)
}

// Responsible for deciding whether to get data from the database or the Codeforces API.
// Chooses database if available for faster response.
type Store struct {
	api API
	db  DB
}

func New(api API, db DB) *Store {
	return &Store{
		api: api,
		db:  db,
	}
}
