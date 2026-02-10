package fetcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/db"
)

func (s *Service) FetchContest(id int) error {
	contestants, contest, err := s.contestProvider.GetContestStandings(context.TODO(), id)
	if err != nil {
		return fmt.Errorf("getting contest %d standings: %w", id, err)
	}

	ratings, err := s.contestProvider.GetContestRatingChanges(context.TODO(), id)
	if err != nil {
		if errors.Is(err, codeforces.ErrRatingChangesUnavailable) {
			err = errors.Join(err, s.insertContestDB(context.Background(), contest, nil))
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
		const minOldTime = 24 * 14 * time.Hour // Two weeks
		isOld := contest.StartTime.Add(minOldTime).Before(time.Now())
		if isOld {
			// Only insert the contest, no contestants as the don't have any rating info.
			// This is to avoid calling the Codeforces API many times for the same contest,
			// when we could just store it to indicate that we already have all available data.
			_, err = s.contestRepo.UpsertContest(context.TODO(), contest)
			return errors.Join(fmt.Errorf(
				"contest %d: %w, but is old so will store in db to avoid future fetches", id, ErrNoRatingInfo),
				err)
		}

		return fmt.Errorf("contest %d: %w", id, ErrNoRatingInfo)
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
	return s.insertContestDB(context.TODO(), contest, contestants)
}

// @param maxAge The maximum age allowed before the data is considered stale.
// @return Slice of the IDs of all contests needing to be updated. (Either stale or not previously fetched).
func (s *Service) FindContestsToUpdate(maxAge time.Duration) ([]int, error) {
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

	stale, err := s.contestRepo.FindStaleContests(context.TODO(), maxAge)
	if err != nil {
		return nil, fmt.Errorf("finding stale contests: %w", err)
	}
	result = append(result, stale...)

	return result, nil
}

func (s *Service) insertContestDB(ctx context.Context, contest *codeforces.Contest,
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
