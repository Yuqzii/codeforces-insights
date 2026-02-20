package handlers

import (
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

	var data struct {
		Count     int `json:"count"`
		MinRating int `json:"minRating"`
		MaxRating int `json:"maxRating"`
		Contests  []struct {
			ID      int      `json:"id"`
			Indices []string `json:"indices"`
		} `json:"contests"`
	}
	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
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

			if errors.Is(err, db.ErrNoProblemsForContest) {
				log.Printf("Couldn't find problems from contest %d: %v\n", c.ID, err)
				continue
			}

			if errors.Is(err, recommender.ErrInvalidIndices) {
				http.Error(w, "Invalid problem indices", http.StatusBadRequest)
				return
			}

			log.Printf("Error finding unsolved problem for contest %d and indices %v: %v\n",
				c.ID, c.Indices, err)
		}
		unsolved = append(unsolved, p)
	}

	if len(unsolved) == 0 {
		http.Error(w, "Failure finding unsolved problems", http.StatusInternalServerError)
		log.Println("Recommend request failed due to inability to find unsolved problems")
		return
	}

	recs, err := h.rec.Recommend(ctx, unsolved, data.Count, data.MinRating, data.MaxRating)
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
