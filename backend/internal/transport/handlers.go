package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type Client interface {
	GetUser(context.Context, string) (*codeforces.User, error)
}

type StatsService interface {
	Categories(string) (map[string]int, error)
	Ratings(string) (map[int]int, error)
}

type Handler struct {
	client Client
	stats  StatsService
}

func NewHandler(api Client, stats StatsService) *Handler {
	return &Handler{
		client: api,
		stats:  stats,
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
	ratings, err := h.stats.Ratings(handle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	j, err := json.Marshal(ratings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
