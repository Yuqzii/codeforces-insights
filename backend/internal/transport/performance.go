package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/yuqzii/cf-stats/internal/stats"
)

type perfManager struct {
	jobs chan perfJob
	mu   sync.Mutex

	crp ContestResultsProvider
}

type perfJob struct {
	ctx context.Context
	chn chan<- perfResult

	contestID int
	rank      int
	rating    int
}

type perfResult struct {
	performance int
	err         error
}

func (h *Handler) HandleGetPerformance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	handle := r.PathValue("handle")
	ratings, err := h.client.GetRatingChanges(ctx, handle)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error getting rating history for performance: %v\n", err)
		}
		return
	}

	type performance struct {
		Rating    int `json:"rating"`
		Timestamp int `json:"timestamp"`
	}

	perf := make([]performance, len(ratings))
	for i := range ratings {
		contestants, contest, err := h.crp.GetContestResults(ctx, ratings[i].ContestID)
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("Error getting contest %d results for performance: %v\n", ratings[i].ContestID, err)
			}
			return
		}

		seed := stats.CalculateSeed(contestants, contest)
		perf[i].Rating = seed.CalculatePerformance(ratings[i].Rank, ratings[i].OldRating)
		perf[i].Timestamp = ratings[i].Timestamp
	}

	j, err := json.Marshal(perf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshalling performance: %v\n", err)
		return
	}

	if _, err = w.Write(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error writing performance: %v\n", err)
		return
	}
}

func (p *perfManager) makeJob(ctx context.Context, id int, rank int, rating int) <-chan perfResult {
	resChan := make(chan perfResult)

	p.jobs <- perfJob{
		ctx:       ctx,
		chn:       resChan,
		contestID: id,
		rank:      rank,
		rating:    rating,
	}

	return resChan
}

func (p *perfManager) perfWorker() {
	for {
		job := <-p.jobs

		select {
		case <-job.ctx.Done():
			continue
		default:
		}

		contestants, contest, err := p.crp.GetContestResults(context.Background(), job.contestID)
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				job.chn <- perfResult{
					err: err,
				}
				close(job.chn)
			}
			return
		}

		seed := stats.CalculateSeed(contestants, contest)
		perf := seed.CalculatePerformance(job.rank, job.rating)

		job.chn <- perfResult{
			performance: perf,
		}
		close(job.chn)
	}
}
