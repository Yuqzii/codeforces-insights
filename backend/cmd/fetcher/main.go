package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yuqzii/codeforces-insights/internal/codeforces"
	"github.com/yuqzii/codeforces-insights/internal/db"
	"github.com/yuqzii/codeforces-insights/internal/fetcher"
)

const (
	dbHost string = "postgres"
	dbPort uint16 = 5432
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
		codeforces.WithIntervals(2*time.Second, 1*time.Minute),
	)

	f := fetcher.New(cfClient, db, cfClient, db, db)

	fetchContests := flag.Bool("contests", false, "Should we fetch contests?")
	fetchProblems := flag.Bool("problems", false, "Should we fetch problems?")
	maxContestsAge := flag.Duration("maxContestAge", time.Since(time.Time{}),
		"Past what age should contests be re-fetched? Default is no maximum.")
	maxContestUpdates := flag.Int("maxContestUpdates", -1,
		"Maximum allowed contests to update. Default is no maximum")
	retryCount := flag.Int("retryCount", 5,
		"How many times should we retry fetching a contest if it gives an error?")
	contestIDToFetch := flag.Int("contestID", -1, "Fetch this specific contest. Used for debugging.")
	flag.Parse()

	if *fetchContests {
		log.Println("Finding unfetched contests")
		contests, err := f.FindContestsToUpdate(context.Background(), *maxContestsAge)
		if err != nil {
			log.Fatalf("Failed to find contests to update: %v\n", err)
		}

		if *maxContestUpdates != -1 && *maxContestUpdates < len(contests) {
			// Limit updates to maxContestUpdates.
			contests = contests[:*maxContestUpdates]
		}

		if *contestIDToFetch != -1 {
			contest, err := f.FindSpecificContest(context.Background(), *contestIDToFetch)
			if err != nil {
				log.Fatalf("Failed to find contest %d: %v", *contestIDToFetch, err)
			}

			contests = contests[:0]
			contests = append(contests, &contest)
		}

		log.Printf("Starting fetching for %d contests\n", len(contests))
		failCnt := 0

		f := fetcher.New(cfClient, db, cfClient, db, db)

		i := 0
		curFail := 0
		for i < len(contests) {
			err := f.FetchContest(context.Background(), contests[i])
			shouldContinue := true
			if err != nil {
				if errors.Is(err, codeforces.ErrRatingChangesUnavailable) {
					// Usually means contest was unrated, we can ignore this.
				} else if errors.Is(err, codeforces.ErrCFServerProblem) {
					curFail++
					if curFail <= *retryCount {
						// Try fetching current contest again.
						shouldContinue = false
						log.Printf("Fetching contest %d failed, retrying: %v\n", contests[i].ID, err)
					} else {
						failCnt++
						log.Printf("Fetching contest %d exceeded retry limit (%d): %v\n",
							contests[i].ID, *retryCount, err)
					}
				} else {
					failCnt++
					log.Printf("Failed to fetch contest %d: %v\n", contests[i].ID, err)
				}
			} else {
				log.Printf("Successfully fetched contest %d (%d/%d)\n", contests[i].ID, i+1, len(contests))
			}

			if shouldContinue {
				i++
				curFail = 0
			}
		}

		outputStr := fmt.Sprintf("Fetched %d/%d contests", len(contests)-failCnt, len(contests))
		log.Println(outputStr)
	}

	if *fetchProblems {
		log.Println("Fetching problems...")

		count, err := f.FetchProblems(context.TODO())
		if err != nil {
			log.Printf("Failed to fetch problems: %s", err)
		} else {
			log.Printf("Successfully fetched and updated %d problems\n", count)
		}
	}
}
