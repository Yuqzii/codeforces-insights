package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/stats"
)

type Client interface {
	GetUser(context.Context, string) (*codeforces.User, error)
	GetSubmissions(context.Context, string) ([]codeforces.Submission, error)
}

type Handler struct {
	client Client
}

func NewHandler(api Client) *Handler {
	return &Handler{
		client: api,
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
	}

	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h *Handler) HandleGetRatings(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	s, err := h.client.GetSubmissions(context.TODO(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	solved := stats.FilterSolved(s)
	ratings := stats.SolvedRatings(solved)

	j, err := json.Marshal(ratings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h *Handler) HandleGetTags(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	s, err := h.client.GetSubmissions(context.TODO(), handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	solved := stats.FilterSolved(s)
	tags := stats.SolvedTags(solved)

	j, err := json.Marshal(tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
