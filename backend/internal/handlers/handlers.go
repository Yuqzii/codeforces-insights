package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/recommender"
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
	GetContestResultsFromHandle(ctx context.Context, handle string) (
		[]codeforces.Contestant, error)
}

type PercentileProvider interface {
	GetPercentile(rating int) float64
}

type RecommendationProvider interface {
	Recommend(ctx context.Context, probs []*codeforces.Problem, cnt, minRat, maxRat int) (
		[]*recommender.ProbWithScore, error)
	FindFirstUnsolvedProblem(ctx context.Context, contestID int, indices []string) (
		*codeforces.Problem, error)
}

type Handler struct {
	client     Client
	crp        ContestResultsProvider
	perf       perfManager
	percentile PercentileProvider
	rec        RecommendationProvider
}

func New(api Client, crp ContestResultsProvider, percentile PercentileProvider, rec RecommendationProvider,
	perfJobsBuffer int, perfWorkers int) *Handler {
	h := &Handler{
		client: api,
		crp:    crp,
		perf: perfManager{
			jobs: make(chan perfJob, perfJobsBuffer),
			crp:  crp,
		},
		percentile: percentile,
		rec:        rec,
	}

	for range perfWorkers {
		go h.perf.worker()
	}

	return h
}

func (h *Handler) HandlePercentile(w http.ResponseWriter, r *http.Request) {
	rating, err := strconv.Atoi(r.PathValue("rating"))
	if err != nil {
		http.Error(w, "rating must be a number", http.StatusBadRequest)
		return
	}

	percentile := h.percentile.GetPercentile(rating)

	j, err := json.Marshal(percentile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error marshalling percentile json: %v", err)
		return
	}

	if _, err = w.Write(j); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error writing percentile: %v", err)
		return
	}
}
