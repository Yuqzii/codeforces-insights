package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/recommender"
)

const maxRecommendRequestSize = 1 << 16 // 65536 bytes

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
		http.Error(w, "Failure reading request", http.StatusInternalServerError)
		log.Printf("Error reading recommend request: %v\n", err)
		return
	}

	var data struct {
		Count    int `json:"count"`
		Contests []struct {
			ID      int      `json:"id"`
			Indices []string `json:"indices"`
		} `json:"contests"`
	}
	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Failure reading request", http.StatusInternalServerError)
		log.Printf("Error unmarshalling json in recommend request: %v\n", err)
		return
	}

	if data.Count > 10 || data.Count < 0 {
		http.Error(w, "Invalid count requested, must be between 0 and 10", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Find unsolved problem for each contest.
	unsolved := make([]*codeforces.Problem, 0, len(data.Contests))
	for _, c := range data.Contests {
		p, err := h.rec.FindFirstUnsolvedProblem(ctx, c.ID, c.Indices)
		if err != nil {
			if errors.Is(err, recommender.ErrNoUnsolvedProblem) {
				continue
			}
			http.Error(w, "Failure finding unsolved problems", http.StatusInternalServerError)
			log.Printf("Error finding unsolved problem for contest %d and indices %v: %v",
				c.ID, c.Indices, err)
		}
		unsolved = append(unsolved, p)
	}

	recs, err := h.rec.Recommend(ctx, unsolved, data.Count)
	if err != nil {
		http.Error(w, "Failure recommending problems", http.StatusInternalServerError)
		log.Printf("Error recommending %d problems: %v\n", data.Count, err)
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
