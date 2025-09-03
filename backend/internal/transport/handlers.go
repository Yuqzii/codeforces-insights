package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type API interface {
	GetUser(context.Context, string) (*codeforces.User, error)
}

type Handler struct {
	api API
}

func NewHandler(api API) *Handler {
	return &Handler{
		api: api,
	}
}

func (h *Handler) HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func (h *Handler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	handle := r.PathValue("handle")
	user, err := h.api.GetUser(context.TODO(), handle)
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
