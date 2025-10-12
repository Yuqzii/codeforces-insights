package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/db"
	"github.com/yuqzii/cf-stats/internal/handlers"
	"github.com/yuqzii/cf-stats/internal/store"
)

const (
	dbHost string = "postgres"
	dbPort uint16 = 5432

	cfTimeBetweenReqs time.Duration = 2 * time.Second

	perfJobsBuffer int = 1000
	perfWorkerCnt  int = 10
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
		cfTimeBetweenReqs,
	)

	store := store.New(cfClient, db)

	log.Println("Setting up API handler")
	h := handlers.New(cfClient, store, perfJobsBuffer, perfWorkerCnt)

	log.Println("Setting up server")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{handle}", h.HandleGetUser)
	mux.HandleFunc("GET /users/solved-ratings/{handle}", h.HandleGetRatings)
	mux.HandleFunc("GET /users/solved-tags/{handle}", h.HandleGetTags)
	mux.HandleFunc("GET /users/rating/{handle}", h.HandleGetRatingChanges)
	mux.HandleFunc("GET /users/performance/{handle}", h.HandleGetPerformance)
	mux.HandleFunc("GET /users/solved-ratings-time/{handle}", h.HandleGetRatingTime)

	log.Println("Server listening on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", mux))
}
