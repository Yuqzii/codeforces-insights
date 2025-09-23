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
	GetContestStandings(ctx context.Context, id int) ([]codeforces.Contestant, *codeforces.Contest, error)
}

type ContestRepository interface {
	UpsertContestTx(context.Context, db.Querier, *codeforces.Contest) (id int, err error)
	InsertContestResultsTx(context.Context, db.Querier, []codeforces.Contestant, int) error
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
