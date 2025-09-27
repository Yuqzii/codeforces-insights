package transport

import (
	"context"
	"encoding/json"
	"fmt"
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
}

func NewHandler(api Client, crp ContestResultsProvider) *Handler {
	return &Handler{
		client: api,
		crp:    crp,
	}
}

func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	user, err := h.client.GetUser(context.TODO(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h *Handler) HandleGetRatings(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	s, err := h.client.GetSubmissions(context.TODO(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	solved := stats.FilterSolved(s)
	ratings := stats.SolvedRatings(solved)

	j, err := json.Marshal(ratings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h *Handler) HandleGetTags(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	s, err := h.client.GetSubmissions(context.TODO(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	solved := stats.FilterSolved(s)
	tags := stats.SolvedTags(solved)

	j, err := json.Marshal(tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h *Handler) HandleGetTagsAndRatings(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	s, err := h.client.GetSubmissions(context.TODO(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	solved := stats.FilterSolved(s)
	tags := stats.SolvedTags(solved)
	ratings := stats.SolvedRatings(solved)

	type tagsAndRatings struct {
		Tags    []stats.Tag `json:"tags"`
		Ratings map[int]int `json:"ratings"`
	}
	combined := tagsAndRatings{
		Tags:    tags,
		Ratings: ratings,
	}
	j, err := json.Marshal(combined)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h *Handler) HandleGetRatingChanges(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	ratings, err := h.client.GetRatingChanges(context.TODO(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(ratings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h *Handler) HandleGetPerformance(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	ratings, err := h.client.GetRatingChanges(context.TODO(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type performance struct {
		Rating    int `json:"rating"`
		Timestamp int `json:"timestamp"`
	}

	perf := make([]performance, len(ratings))
	for i := range ratings {
		contestants, contest, err := h.crp.GetContestResults(
			context.TODO(), ratings[i].ContestID,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		seed := stats.CalculateSeed(contestants, contest)
		perf[i].Rating = seed.CalculatePerformance(ratings[i].Rank, ratings[i].OldRating)
		perf[i].Timestamp = ratings[i].Timestamp
	}

	j, err := json.Marshal(perf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h *Handler) HandleGetRatingTime(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	s, err := h.client.GetSubmissions(context.TODO(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
