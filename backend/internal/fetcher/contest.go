package fetcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yuqzii/codeforces-insights/internal/codeforces"
	"github.com/yuqzii/codeforces-insights/internal/db"
)

var ErrContestNotInList = errors.New("contest is not in Codeforces's list")

func (f *fetcher) FetchContest(ctx context.Context, contest *codeforces.Contest) error {
	contestants, _, err := f.contestProvider.GetContestStandings(ctx, contest.ID)
	if err != nil {
		return fmt.Errorf("getting contest standings: %w", err)
	}

	ratings, err := f.contestProvider.GetContestRatingChanges(ctx, contest.ID)
	if err != nil {
		if errors.Is(err, codeforces.ErrRatingChangesUnavailable) {
			// Insert to avoid refetch, usually means the contest was unrated.
			_, err = f.contestRepo.UpsertContest(ctx, contest)
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
		if contestIsOld(contest) {
			// Only insert the contest, no contestants as they don't have any rating info.
			// This is to avoid calling the Codeforces API many times for the same contest,
			// when we could just store it to indicate that we already have all available data.
			_, err = f.contestRepo.UpsertContest(ctx, contest)
			return err
		}

		return fmt.Errorf("contest %d: %w", contest.ID, ErrNoRatingInfo)
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

	return f.insertContestDB(ctx, contest, contestants)
}

// @param maxAge The maximum age allowed before the data is considered stale.
// @return Slice of the IDs of all contests needing to be updated. (Either stale or not previously fetched).
func (f *fetcher) FindContestsToUpdate(ctx context.Context, maxAge time.Duration) ([]*codeforces.Contest, error) {
	contests, err := f.contestProvider.GetContests(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting contests: %w", err)
	}

	finished := make(map[int]*codeforces.Contest, 0)
	for i := range contests {
		if contests[i].Phase == "FINISHED" {
			finished[contests[i].ID] = &contests[i]
		}
	}

	finishedIDs := make([]int, 0, len(finished))
	for id := range finished {
		finishedIDs = append(finishedIDs, id)
	}

	existing, err := f.contestRepo.ContestsExists(ctx, finishedIDs)
	if err != nil {
		return nil, fmt.Errorf("checking contests existence: %w", err)
	}

	result := make([]*codeforces.Contest, 0)
	for id, c := range finished {
		_, exists := existing[id]
		if !exists {
			result = append(result, c)
		}
	}

	stale, err := f.contestRepo.FindStaleContests(ctx, maxAge)
	if err != nil {
		return nil, fmt.Errorf("finding stale contests: %w", err)
	}
	for _, id := range stale {
		result = append(result, finished[id])
	}

	return result, nil
}

func (f *fetcher) FindSpecificContest(ctx context.Context, id int) (codeforces.Contest, error) {
	contests, err := f.contestProvider.GetContests(ctx)
	if err != nil {
		return codeforces.Contest{}, err
	}

	for _, c := range contests {
		if c.ID == id {
			return c, nil
		}
	}

	return codeforces.Contest{}, ErrContestNotInList
}

// Upserts the provided contest and contestants.
// Uses a gouroutine to avoid blocking while waiting for the DB update.
func (f *fetcher) insertContestDB(ctx context.Context, contest *codeforces.Contest,
	contestants []codeforces.Contestant) error {

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
		return fmt.Errorf("updating db during contest fetch (id %d): %w", contest.ID, err)
	}

	return nil
}

func contestIsOld(contest *codeforces.Contest) bool {
	const minOldTime = 24 * 14 * time.Hour // Two weeks
	isOld := contest.StartTime.Add(minOldTime).Before(time.Now())
	return isOld
}
