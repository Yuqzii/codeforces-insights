package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/db"
)

const (
	dbHost string = "postgres"
	dbPort uint16 = 5432

	cfRequestsPerSecond float64 = 0.4
	cfMaxBurst          int     = 1
)

func main() {
	log.Println("Connecting to database...")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPswd := os.Getenv("POSTGRES_PASSWORD")
	db, err := db.New(context.Background(), dbHost, dbUser, dbPswd, dbName, dbPort)
	if err != nil {
		log.Fatalf("Could not connect to database: %v\n", err)
	}
	log.Println("Connected to database")
	defer db.Close()

	cfClient := codeforces.NewClient(
		http.DefaultClient,
		"https://codeforces.com/api/",
		cfRequestsPerSecond,
		cfMaxBurst,
	)

	fetcher := newFetcher(cfClient, db, db)
	log.Println("Finding unfetched contests")
	unfetched, err := fetcher.findUnfetchedContests()
	if err != nil {
		log.Fatalf("Failed to find unfetched contests: %v\n", err)
	}

	log.Printf("Starting fetching for %d contests\n", len(unfetched))
	bar := progressbar.Default(int64(len(unfetched)), "Fetching contests...")
	failCnt, noRatingCnt := 0, 0
	for _, id := range unfetched {
		err = fetcher.fetchContest(id)
		bar.Add(1)
		if err != nil {
			if errors.Is(err, codeforces.ErrRatingChangesUnavailable) {
				// Usually means contest was unrated
				noRatingCnt++
				continue
			}
			failCnt++
			log.Printf("Failed to fetch contest %d: %v\n", id, err)
		}
	}

	outputStr := fmt.Sprintf("Fetched %d/%d contests", len(unfetched)-failCnt, len(unfetched))
	if noRatingCnt > 0 {
		outputStr += fmt.Sprintf(" (%d was not stored due to missing rating data)", noRatingCnt)
	}
	log.Println(outputStr)
}
