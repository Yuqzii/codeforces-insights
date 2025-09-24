package main

import (
	"context"
	"fmt"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/db"
)

type fetcher struct {
	cp          ContestProvider
	contestRepo ContestRepository
	tx          db.TxManager
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

func NewFetcher(cp ContestProvider, contestRepo ContestRepository, tx db.TxManager) *fetcher {
	return &fetcher{
		cp:          cp,
		contestRepo: contestRepo,
		tx:          tx,
	}
}

func (f *fetcher) fetchContest(id int) error {
	contestants, contest, err := f.cp.GetContestStandings(context.TODO(), id)
	if err != nil {
		return fmt.Errorf("getting contest standings: %w", err)
	}
	rankMap := make(map[int]*codeforces.Contestant)
	for i := range contestants {
		rankMap[contestants[i].Rank] = &contestants[i]
	}

	// Set ratings of contestants
	ratings, err := f.cp.GetContestRatingChanges(context.TODO(), id)
	if err != nil {
		return fmt.Errorf("getting contest ratings: %w", err)
	}
	for _, r := range ratings {
		rankMap[r.Rank].NewRating = r.NewRating
		rankMap[r.Rank].OldRating = r.OldRating
	}

	// Insert to DB in a transaction
	err = f.tx.WithTx(context.TODO(), func(q db.Querier) error {
		id, err := f.contestRepo.UpsertContestTx(context.TODO(), q, contest)
		if err != nil {
			return fmt.Errorf("upserting contest: %w", err)
		}

		err = f.contestRepo.InsertContestResultsTx(context.TODO(), q, contestants, id)
		if err != nil {
			return fmt.Errorf("inserting contest results: %w", err)
		}

		return nil
	})

	return err
}

// @return Slice of the IDs of all unfetched contests.
func (f *fetcher) findUnfetchedContests() ([]int, error) {
	c, err := f.cp.GetContests(context.TODO())
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

	return result, nil
}
