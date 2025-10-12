package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/stats"
)

type perfManager struct {
	jobs chan perfJob
	crp  ContestResultsProvider
}

type perfJob struct {
	ctx context.Context
	chn chan<- perfResult

	contestID int
	rank      int
	rating    int
	timestamp int
}

type perfResult struct {
	performance int
	timestamp   int
	err         error
}

func (h *Handler) HandleGetPerformance(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	handle := r.PathValue("handle")
	ratings, err := h.client.GetRatingChanges(ctx, handle)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error getting rating history for performance: %v\n", err)
		}
		cancel()
		return
	}

	type performance struct {
		Rating    int `json:"rating"`
		Timestamp int `json:"timestamp"`
	}

	perf := make([]performance, 0, len(ratings))
	resChan := make(chan perfResult, len(ratings))

	for i := range ratings {
		h.perf.addJob(ctx, &ratings[i], resChan)
	}

	for range ratings {
		select {
		case perfRes := <-resChan:
			if perfRes.err != nil {
				http.Error(w, perfRes.err.Error(), http.StatusInternalServerError)
				log.Printf("Error getting performance: %v\n", perfRes.err)
				// Cancel context so we don't make unnecessary calculations, and avoids leaking channel
				cancel()
				close(resChan)
				return
			}
			perf = append(perf, performance{
				Rating:    perfRes.performance,
				Timestamp: perfRes.timestamp,
			})
		case <-ctx.Done():
			cancel()
			return
		}
	}

	j, err := json.Marshal(perf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshalling performance: %v\n", err)
		cancel()
		return
	}

	if _, err = w.Write(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error writing performance: %v\n", err)
		cancel()
		return
	}

	cancel()
}

func (p *perfManager) addJob(ctx context.Context, r *codeforces.RatingChange, resChan chan<- perfResult) {
	p.jobs <- perfJob{
		ctx:       ctx,
		chn:       resChan,
		contestID: r.ContestID,
		rank:      r.Rank,
		rating:    r.OldRating,
		timestamp: r.Timestamp,
	}
}

func (p *perfManager) worker() {
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
			}
			return
		}

		seed := stats.CalculateSeed(contestants, contest)
		perf := seed.CalculatePerformance(job.rank, job.rating)

		job.chn <- perfResult{
			performance: perf,
			timestamp:   job.timestamp,
		}
	}
}
