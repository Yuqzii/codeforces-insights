package fetcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/yuqzii/codeforces-insights/internal/codeforces"
	"github.com/yuqzii/codeforces-insights/internal/db"
)

func (f *fetcher) FetchContest(id int) error {
	contestants, contest, err := f.contestProvider.GetContestStandings(context.TODO(), id)
	if err != nil {
		return fmt.Errorf("getting contest standings: %w", err)
	}

	ratings, err := f.contestProvider.GetContestRatingChanges(context.TODO(), id)
	if err != nil {
		if errors.Is(err, codeforces.ErrRatingChangesUnavailable) {
			// Insert to avoid refetch.
			f.insertContestDB(context.Background(), contest, nil)
			return err
		}
		return fmt.Errorf("getting contest ratings: %w", err)
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
			// Only insert the contest, no contestants as they don't have any rating info.
			// This is to avoid calling the Codeforces API many times for the same contest,
			// when we could just store it to indicate that we already have all available data.
			_, err = f.contestRepo.UpsertContest(context.TODO(), contest)
			return err
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

	f.insertContestDB(context.TODO(), contest, contestants)

	return nil
}

// @param maxAge The maximum age allowed before the data is considered stale.
// @return Slice of the IDs of all contests needing to be updated. (Either stale or not previously fetched).
func (f *fetcher) FindContestsToUpdate(maxAge time.Duration) ([]int, error) {
	c, err := f.contestProvider.GetContests(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("getting contests: %w", err)
	}

	finished := make([]int, 0)
	for _, cont := range c {
		if cont.Phase == "FINISHED" && !containsCyrillic(cont.Name) {
			finished = append(finished, cont.ID)
		}
	}

	existing, err := f.contestRepo.ContestsExists(context.TODO(), finished)
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

	stale, err := f.contestRepo.FindStaleContests(context.TODO(), maxAge)
	if err != nil {
		return nil, fmt.Errorf("finding stale contests: %w", err)
	}
	result = append(result, stale...)

	return result, nil
}

// Upserts the provided contest and contestants.
// Uses a gouroutine to avoid blocking while waiting for the DB update.
func (f *fetcher) insertContestDB(ctx context.Context, contest *codeforces.Contest,
	contestants []codeforces.Contestant) {

	go func() {
		err := f.tx.WithTx(ctx, func(q db.Querier) error {
			id, err := f.contestRepo.UpsertContestTx(ctx, q, contest)
			if err != nil {
				return fmt.Errorf("upserting contest %d: %w", id, err)
			}

			err = f.contestRepo.InsertContestResultsTx(ctx, q, contestants, id)
			if err != nil {
				return fmt.Errorf("inserting contest %d results: %w", id, err)
			}

			return nil
		})

		if err != nil {
			log.Printf("Error when updating db during contest fetch (id %d): %v\n", contest.ID, err)
		} else {
			log.Printf("Successfully updated contest %d\n", contest.ID)
		}
	}()
}
