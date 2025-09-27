package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/db"
	"github.com/yuqzii/cf-stats/internal/store"
	"github.com/yuqzii/cf-stats/internal/transport"
)

const (
	dbHost string = "postgres"
	dbPort uint16 = 5432

	cfRequestsPerSecond float64 = 0.5
	cfMaxBurst          int     = 1
)

func main() {
	log.Println("Connecting to database...")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPswd := os.Getenv("POSTGRES_PASSWORD")
	db, err := db.New(context.Background(), dbHost, dbUser, dbPswd, dbName, dbPort)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	log.Println("Connected to database")
	defer db.Close()

	log.Println("Setting up Codeforces client")
	cfClient := codeforces.NewClient(
		http.DefaultClient,
		"https://codeforces.com/api/",
		cfRequestsPerSecond,
		cfMaxBurst)

	store := store.New(cfClient, db)

	log.Println("Setting up API handler")
	h := transport.NewHandler(cfClient, store)

	log.Println("Setting up server")
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.HandleRoot)
	mux.HandleFunc("GET /users/{handle}", h.HandleGetUser)
	mux.HandleFunc("GET /users/solved-ratings/{handle}", h.HandleGetRatings)
	mux.HandleFunc("GET /users/solved-tags/{handle}", h.HandleGetTags)
	// Prefer calling this to minimize Codeforces API calls.
	mux.HandleFunc("GET /users/solved-tags-ratings/{handle}", h.HandleGetTagsAndRatings)
	mux.HandleFunc("GET /users/rating/{handle}", h.HandleGetRatingChanges)
	mux.HandleFunc("GET /users/performance/{handle}", h.HandleGetPerformance)
	mux.HandleFunc("GET /users/solved-ratings-time/{handle}", h.HandleGetRatingTime)

	log.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", mux)
}
