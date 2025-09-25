package fetcher

import (
	"sync"

	"github.com/yuqzii/cf-stats/internal/db"
)

func worker(cp ContestProvider, cr ContestRepository, tx db.TxManager, jobs <-chan int, results chan<- error) {
	f := New(cp, cr, tx)
	for j := range jobs {
		err := f.FetchContest(j)
		results <- err
	}
}

func CreateWorkers(cnt int, ids []int, cp ContestProvider, cr ContestRepository, tx db.TxManager) <-chan error {
	jobs := make(chan int)
	results := make(chan error)

	var wg sync.WaitGroup

	for range cnt {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(cp, cr, tx, jobs, results)
		}()
	}

	go func() {
		for _, id := range ids {
			jobs <- id
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
