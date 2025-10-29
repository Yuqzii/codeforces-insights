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
	"github.com/yuqzii/cf-stats/internal/stats"
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

	log.Println("Calculating percentiles")
	cfUsers, err := cfClient.GetActiveUsers(context.Background())
	if err != nil {
		log.Fatalf("Could not get active Codeforces users: %v\n", err)
	}
	percentile := stats.NewPercentile(cfUsers)

	log.Println("Setting up API handler")

	h := handlers.New(cfClient, store, percentile, perfJobsBuffer, perfWorkerCnt)

	log.Println("Setting up server")
	mux := http.NewServeMux()
	mux.HandleFunc("POST /performance", h.HandlePerformance)
	mux.HandleFunc("GET /percentile/{rating}", h.HandlePercentile)

	log.Println("Server listening on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", mux))
}
