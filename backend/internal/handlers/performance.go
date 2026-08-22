package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/yuqzii/codeforces-insights/internal/codeforces"
	"github.com/yuqzii/codeforces-insights/internal/stats"
)

const maxPerfRequestSize = 1 << 16 // 65536 bytes

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

func (h *Handler) HandlePerformance(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxPerfRequestSize)
	defer r.Body.Close() //nolint:errcheck
	body, err := io.ReadAll(r.Body)
	if err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error reading performance request: %v\n", err)
		return
	}

	var data struct {
		Handle  string                    `json:"handle"`
		Ratings []codeforces.RatingChange `json:"ratingHistory"`
	}
	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithCancel(r.Context())

	type performance struct {
		Rating    int `json:"rating"`
		Timestamp int `json:"timestamp"`
	}

	idToRat := make(map[int]*codeforces.RatingChange)
	for i := range data.Ratings {
		idToRat[data.Ratings[i].ContestID] = &data.Ratings[i]
	}

	results, err := h.crp.GetContestResultsFromHandle(ctx, data.Handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error getting contest results from handle '%s': %v", data.Handle, err)
		cancel()
		return
	}

	// Update ratings with the correct rank that we store.
	for _, r := range results {
		rating, ok := idToRat[r.ContestID]
		if ok {
			rating.Rank = r.Rank
		}
	}

	perf := make([]performance, 0, len(data.Ratings))
	resChan := make(chan perfResult, len(data.Ratings))
	for i := range data.Ratings {
		h.perf.addJob(ctx, &data.Ratings[i], resChan)
	}

	for range data.Ratings {
		select {
		case perfRes := <-resChan:
			if perfRes.err != nil {
				http.Error(w, perfRes.err.Error(), http.StatusInternalServerError)
				log.Printf("Error getting performance: %v\n", perfRes.err)
				// Cancel context so we don't make unnecessary calculations.
				cancel()
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

		contestants, contest, err := p.crp.GetContestResults(job.ctx, job.contestID)
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

		select {
		case job.chn <- perfResult{
			performance: perf,
			timestamp:   job.timestamp,
		}:
		case <-job.ctx.Done():
		}
	}
}
