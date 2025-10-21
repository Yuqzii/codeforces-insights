package fetcher

import (
	"context"
	"errors"
	"fmt"

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

func (s *Service) FetchContest(id int) error {
	contestants, contest, err := s.contestProvider.GetContestStandings(context.TODO(), id)
	if err != nil {
		return fmt.Errorf("getting contest %d standings: %w", id, err)
	}

	ratings, err := s.contestProvider.GetContestRatingChanges(context.TODO(), id)
	if err != nil {
		if errors.Is(err, codeforces.ErrRatingChangesUnavailable) {
			err = errors.Join(err, s.insertDB(context.Background(), contest, nil))
			return err
		}
		return fmt.Errorf("getting contest %d ratings: %w", id, err)
	}

	hasRatingInfo := false
	ratingMap := make(map[string]*codeforces.RatingChange)
	for i := range ratings {
		if ratings[i].OldRating != 0 {
			hasRatingInfo = true
		}
		ratingMap[ratings[i].Handle] = &ratings[i]
	}

	if !hasRatingInfo {
		return ErrNoRatingInfo
	}

	// Set ratings of contestants
	for i, c := range contestants {
		for _, handle := range c.MemberHandles {
			r, ok := ratingMap[handle]
			// Use rating of party member with maximum previous rating
			if ok && r.OldRating > contestants[i].OldRating {
				contestants[i].OldRating = r.OldRating
				contestants[i].NewRating = r.NewRating
			}
		}
	}

	// Insert to DB in a transaction
	return s.insertDB(context.TODO(), contest, contestants)
}

// @return Slice of the IDs of all unfetched contests.
func (s *Service) FindUnfetchedContests() ([]int, error) {
	c, err := s.contestProvider.GetContests(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("getting contests: %w", err)
	}

	finished := make([]int, 0)
	for _, cont := range c {
		if cont.Phase == "FINISHED" && !containsCyrillic(cont.Name) {
			finished = append(finished, cont.ID)
		}
	}

	existing, err := s.contestRepo.ContestsExists(context.TODO(), finished)
	if err != nil {
		return nil, fmt.Errorf("checking contests existence: %w", err)
	}

	result := make([]int, 0)
	for _, id := range finished {
		_, exists := existing[id]
		if !exists {
			result = append(result, id)
		}
	}

	return result, nil
}

func (s *Service) insertDB(ctx context.Context, contest *codeforces.Contest,
	contestants []codeforces.Contestant) error {

	return s.tx.WithTx(ctx, func(q db.Querier) error {
		id, err := s.contestRepo.UpsertContestTx(ctx, q, contest)
		if err != nil {
			return fmt.Errorf("upserting contest %d: %w", id, err)
		}

		err = s.contestRepo.InsertContestResultsTx(ctx, q, contestants, id)
		if err != nil {
			return fmt.Errorf("inserting contest %d results: %w", id, err)
		}

		return nil
	})
}
