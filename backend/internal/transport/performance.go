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
	jobs      chan int
	listeners map[int][]perfListener
	mu        sync.Mutex

	crp ContestResultsProvider
}

type perfListener struct {
	ctx context.Context
	chn chan<- perfResult
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

func (p *perfManager) makeJob(ctx context.Context, id int) <-chan perfResult {
	resChan := make(chan perfResult)

	_, queued := p.listeners[id]
	if !queued {
		p.listeners[id] = make([]perfListener, 0, 1)
	}

	p.listeners[id] = append(p.listeners[id], perfListener{
		ctx: ctx,
		chn: resChan,
	})

	return resChan
}

func (p *perfManager) perfWorker() {
	for {
		id := <-p.jobs

		if p.listenersCancelled(id) {
			continue
		}

		contestants, contest, err := p.crp.GetContestResults(context.Background(), id)
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				p.sendErrToListeners(err, id)
			}
			return
		}

		seed := stats.CalculateSeed(contestants, contest)
		_ = seed
		//perf := seed.CalculatePerformance(ratings[i].Rank, ratings[i].OldRating)
	}
}

func (p *perfManager) listenersCancelled(id int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, l := range p.listeners[id] {
		select {
		case <-l.ctx.Done():
			continue
		default:
			return false
		}
	}
	return true
}

func (p *perfManager) sendErrToListeners(err error, id int) {
	result := perfResult{
		performance: -1,
		err:         err,
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, l := range p.listeners[id] {
		select {
		case <-l.ctx.Done(): // Don't send to cancelled receiver
		default:
			l.chn <- result
		}
		close(l.chn)
	}
}
