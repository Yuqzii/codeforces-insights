package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/db"
	"github.com/yuqzii/cf-stats/internal/recommender"
)

const maxRecommendRequestSize = 1 << 16 // 65536 bytes

type RecommendationProvider interface {
	Recommend(ctx context.Context, probs []*codeforces.Problem, disallowedProbs map[int64]struct{},
		cnt, minRat, maxRat int) ([]*recommender.ProbWithScore, error)
	FindFirstUnsolvedProblem(ctx context.Context, contestID int, indices []string) (
		*codeforces.Problem, error)
	FindSolvedRecentContests(subs []codeforces.Submission, lookback int) map[int][]*codeforces.Problem
}

type recommendReq struct {
	Count        int                     `json:"count"`
	MinRating    int                     `json:"minRating"`
	MaxRating    int                     `json:"maxRating"`
	AcceptedSubs []codeforces.Submission `json:"submissions"`
	Lookback     int                     `json:"lookback"`
}

// The endpoint expects json on the format specified in the data struct.
func (h *Handler) HandleRecommend(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxRecommendRequestSize)
	defer r.Body.Close() //nolint:errcheck
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "Couldn't read request", http.StatusBadRequest)
		log.Printf("Error reading recommend request: %v\n", err)
		return
	}

	var data recommendReq
	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Printf("Error unmarshalling json in recommend request: %v\n", err)
		return
	}

	if err = data.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx := r.Context()

	solvedRecent := h.rec.FindSolvedRecentContests(data.AcceptedSubs, data.Lookback)

	// Find first unsolved problem for each from solvedRecent.
	unsolved := make([]*codeforces.Problem, 0, data.Lookback)
	for contestID, probs := range solvedRecent {
		indices := make([]string, 0)
		for _, p := range probs {
			indices = append(indices, p.Index)
		}

		unsolvedProb, err := h.rec.FindFirstUnsolvedProblem(ctx, contestID, indices)
		if err != nil {
			if errors.Is(err, recommender.ErrNoUnsolvedProblem) {
				continue
			}

			if errors.Is(err, db.ErrNoProblemsForContest) {
				log.Printf("Couldn't find problems from contest %d: %v\n", contestID, err)
				continue
			}

			// TODO: Update to more correct error, the client no longer provides the indices.
			if errors.Is(err, recommender.ErrInvalidIndices) {
				http.Error(w, "Invalid problem indices", http.StatusBadRequest)
				return
			}

			log.Printf("Error finding unsolved problem for contest %d and indices %v: %v\n",
				contestID, indices, err)
		}

		unsolved = append(unsolved, unsolvedProb)
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

	if _, err = w.Write(j); err != nil {
		http.Error(w, "Failure writing response", http.StatusInternalServerError)
		log.Printf("Error writing recommendations: %v\n", err)
	}
}

func (r *recommendReq) validate() error {
	if r.Count > 10 || r.Count < 0 {
		return errors.New("count must be between 0 and 10")
	}

	return nil
}
