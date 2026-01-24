package fetcher

import (
	"context"
	"errors"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/db"
)

var ErrNoRatingInfo = errors.New("no rating info exists for this contest")

type Service struct {
	contestProvider ContestProvider
	contestRepo     ContestRepository
	tx              db.TxManager
}

type ContestProvider interface {
	GetContests(context.Context) ([]codeforces.Contest, error)
	GetContestStandings(ctx context.Context, id int) ([]codeforces.Contestant, *codeforces.Contest, error)
	GetContestRatingChanges(ctx context.Context, id int) ([]codeforces.RatingChange, error)
}

type ContestRepository interface {
	UpsertContestTx(context.Context, db.Querier, *codeforces.Contest) (id int, err error)
	UpsertContest(context.Context, *codeforces.Contest) (id int, err error)
	InsertContestResultsTx(context.Context, db.Querier, []codeforces.Contestant, int) error
	ContestsExists(context.Context, []int) (existingIDs map[int]struct{}, err error)
}

func New(cp ContestProvider, contestRepo ContestRepository, tx db.TxManager) *Service {
	return &Service{
		contestProvider: cp,
		contestRepo:     contestRepo,
		tx:              tx,
	}
}
