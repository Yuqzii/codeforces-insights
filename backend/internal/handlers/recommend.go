package handlers

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/recommender"
)

type RecommendationProvider interface {
	Recommend(ctx context.Context, probs []*codeforces.Problem, disallowedProbs map[int64]struct{},
		cnt, minRat, maxRat int) ([]*recommender.ProbWithScore, error)
	FindUnsolvedProblems(ctx context.Context, solvedByContest map[int][]*codeforces.Problem) (
		[]*codeforces.Problem, error)
	FindSolvedRecentContests(subs []codeforces.Submission, lookback int) map[int][]*codeforces.Problem
}

type recommendReq struct {
	Count        int                     `json:"count"`
	MinRating    int                     `json:"minRating"`
	MaxRating    int                     `json:"maxRating"`
	Lookback     int                     `json:"lookback"`
	AcceptedSubs []codeforces.Submission `json:"submissions"`
}

// The endpoint expects json on the format specified in the data struct.
func (h *Handler) HandleRecommend(w http.ResponseWriter, r *http.Request) {
	const (
		maxRequestSize = 1 << 20 // 1 MiB
		maxRequestBody = 1 << 23 // 8 MiB
	)

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	defer r.Body.Close() //nolint:errcheck

	var reader io.Reader = r.Body

	if r.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Invalid gzip", http.StatusBadRequest)
			return
		}
		defer gz.Close() // nolint:errcheck

		reader = http.MaxBytesReader(w, gz, maxRequestBody)
	}

	var data recommendReq
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		http.Error(w, "Couldn't read request", http.StatusBadRequest)
		log.Printf("Error reading recommend request: %v\n", err)
		return
	}

	if err := data.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	solvedByContest := h.rec.FindSolvedRecentContests(data.AcceptedSubs, data.Lookback)
	unsolved, err := h.rec.FindUnsolvedProblems(ctx, solvedByContest)
	if err != nil {
		if errors.Is(err, recommender.ErrInvalidIndices) {
			http.Error(w, "Invalid indices in recent solved problems", http.StatusBadRequest)
			return
		}

		http.Error(w, "Could not find unsolved problems", http.StatusInternalServerError)
		log.Printf("Error in recommend request: %v", err)
		return
	}

	if len(unsolved) == 0 {
		http.Error(w, "Failure finding unsolved problems", http.StatusInternalServerError)
		log.Println("Recommend request failed due to inability to find unsolved problems")
		return
	}

	// Collect all already solved problems into a set.
	disallowedProbs := make(map[int64]struct{})
	for _, sub := range data.AcceptedSubs {
		hash := sub.Problem.Hash()
		disallowedProbs[hash] = struct{}{}
	}

	recs, err := h.rec.Recommend(ctx, unsolved, disallowedProbs, data.Count, data.MinRating, data.MaxRating)
	if err != nil {
		http.Error(w, "Failure recommending problems", http.StatusInternalServerError)
		log.Printf("Error recommending %d problems: %v\n", data.Count, err)
		return
	}

	j, err := json.Marshal(recs)
	if err != nil {
		http.Error(w, "Failure creating response", http.StatusInternalServerError)
		log.Printf("Error marshalling recommendations: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(j); err != nil {
		log.Printf("Error writing recommendations: %v\n", err)
	}
}

func (r *recommendReq) validate() error {
	if r.Count > 10 || r.Count < 0 {
		return errors.New("count must be between 0 and 10")
	}

	if len(r.AcceptedSubs) == 0 {
		return errors.New("submissions cannot be empty")
	}

	return nil
}
