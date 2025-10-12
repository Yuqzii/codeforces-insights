package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/stats"
)

type Client interface {
	GetUser(context.Context, string) (*codeforces.User, error)
	GetSubmissions(context.Context, string) ([]codeforces.Submission, error)
	GetRatingChanges(context.Context, string) ([]codeforces.RatingChange, error)
	GetContestRatingChanges(context.Context, int) ([]codeforces.RatingChange, error)
	GetContestStandings(context.Context, int) ([]codeforces.Contestant, *codeforces.Contest, error)
}

type ContestResultsProvider interface {
	GetContestResults(ctx context.Context, id int) (
		[]codeforces.Contestant, *codeforces.Contest, error)
}

type Handler struct {
	client Client
	crp    ContestResultsProvider
	perf   perfManager
}

func New(api Client, crp ContestResultsProvider, perfJobsBuffer int, perfWorkers int) *Handler {
	h := &Handler{
		client: api,
		crp:    crp,
		perf: perfManager{
			jobs: make(chan perfJob, perfJobsBuffer),
			crp:  crp,
		},
	}

	for range perfWorkers {
		go h.perf.worker()
	}

	return h
}

func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	user, err := h.client.GetUser(r.Context(), handle)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error getting user info: %v\n", err)
		}
		return
	}

	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshalling user info: %v\n", err)
		return
	}

	if _, err = w.Write(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error writing user info: %v\n", err)
		return
	}
}

func (h *Handler) HandleGetRatings(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	s, err := h.client.GetSubmissions(r.Context(), handle)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error getting submissions for ratings: %v\n", err)
		}
		return
	}

	solved := stats.FilterSolved(s)
	ratings := stats.SolvedRatings(solved)

	j, err := json.Marshal(ratings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshalling user solved ratings: %v\n", err)
		return
	}

	if _, err = w.Write(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error writing user solved ratings: %v\n", err)
		return
	}
}

func (h *Handler) HandleGetTags(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	s, err := h.client.GetSubmissions(r.Context(), handle)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error getting submissions for solved tags: %v\n", err)
		}
		return
	}

	solved := stats.FilterSolved(s)
	tags := stats.SolvedTags(solved)

	j, err := json.Marshal(tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshalling user solved tags: %v\n", err)
		return
	}

	if _, err = w.Write(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error writing user solved tags: %v\n", err)
		return
	}
}

func (h *Handler) HandleGetRatingChanges(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	ratings, err := h.client.GetRatingChanges(r.Context(), handle)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error getting user rating history: %v\n", err)
		}
		return
	}

	j, err := json.Marshal(ratings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshalling user rating history: %v\n", err)
		return
	}

	if _, err = w.Write(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error writing user rating history: %v\n", err)
		return
	}
}

func (h *Handler) HandleGetRatingTime(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	s, err := h.client.GetSubmissions(r.Context(), handle)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error getting submissions for solved rating time: %v\n", err)
		}
		return
	}

	solved := stats.FilterSolved(s)
	// Sort by solved time
	slices.SortFunc(solved, func(a, b codeforces.Submission) int {
		return a.Timestamp - b.Timestamp
	})

	type response struct {
		Rating    int `json:"rating"`
		Timestamp int `json:"timestamp"`
	}
	resp := make([]response, 0, len(solved))
	for _, sub := range solved {
		if sub.Problem.Rating == 0 {
			continue
		}
		resp = append(resp, response{
			Rating:    sub.Problem.Rating,
			Timestamp: sub.Timestamp,
		})
	}

	j, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshalling solved ratings time: %v\n", err)
		return
	}

	if _, err = w.Write(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error writing solved ratings time: %v\n", err)
		return
	}
}
