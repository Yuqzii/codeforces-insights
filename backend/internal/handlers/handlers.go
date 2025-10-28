package handlers

import (
	"context"

	"github.com/yuqzii/cf-stats/internal/codeforces"
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
