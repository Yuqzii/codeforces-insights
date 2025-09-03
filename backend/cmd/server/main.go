package main

import (
	"fmt"
	"net/http"

	"github.com/yuqzii/cf-stats/internal/codeforces"
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

	h := transport.NewHandler(cfClient)

	mux.HandleFunc("/", h.HandleRoot)
	mux.HandleFunc("GET /users/{handle}", h.HandleGetUser)

	fmt.Println("Server listening on port :8080")
	http.ListenAndServe(":8080", mux)
}
