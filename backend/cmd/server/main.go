package main

import (
	"fmt"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/stats"
	"github.com/yuqzii/cf-stats/internal/transport"
)

const (
	cfRequestsPerSecond float64 = 0.5
	cfMaxBurst          int     = 1
)

func main() {
	mux := http.NewServeMux()

	cfClient := codeforces.NewClient(
		http.DefaultClient,
		"https://codeforces.com/api/",
		cfRequestsPerSecond,
		cfMaxBurst)

	s := stats.NewService(cfClient)

	h := transport.NewHandler(cfClient, s)

	mux.HandleFunc("/", h.HandleRoot)
	mux.HandleFunc("GET /users/{handle}", h.HandleGetUser)
	mux.HandleFunc("GET /users/solved-ratings/{handle}", h.HandleGetRatings)

	fmt.Println("Server listening on port :8080")
	http.ListenAndServe(":8080", mux)
}
