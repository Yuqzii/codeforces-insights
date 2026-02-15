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

	"github.com/schollz/progressbar/v3"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/db"
	"github.com/yuqzii/cf-stats/internal/fetcher"
)

const (
	dbHost string = "postgres"
	dbPort uint16 = 5432

	workerCnt int = 2
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
		codeforces.WithIntervals(2*time.Second, 20*time.Second),
	)

	f := fetcher.New(cfClient, db, cfClient, db, db)

	fetchContests := flag.Bool("contests", false, "Should we fetch contests?")
	fetchProblems := flag.Bool("problems", false, "Should we fetch problems?")
	maxContestsAge := flag.Duration("maxContestAge", time.Since(time.Time{}),
		"Past what age should contests be re-fetched? Default is no maximum.")
	maxContestUpdates := flag.Int("maxContestUpdates", -1,
		"Maximum allowed contests to update. Default is no maximum")
	flag.Parse()

	if *fetchContests {
		log.Println("Finding unfetched contests")
		contestIDs, err := f.FindContestsToUpdate(*maxContestsAge)
		if err != nil {
			log.Fatalf("Failed to find contests to update: %v\n", err)
		}

		if *maxContestUpdates != -1 && *maxContestUpdates < len(contestIDs) {
			// Limit updates to maxContestUpdates
			contestIDs = contestIDs[:*maxContestUpdates]
		}

		log.Printf("Starting fetching for %d contests\n", len(contestIDs))
		bar := progressbar.Default(int64(len(contestIDs)), "Fetching contests")
		failCnt := 0

		results := fetcher.CreateWorkers(workerCnt, contestIDs, cfClient, db, cfClient, db, db)
		for err := range results {
			bar.Add(1) //nolint:errcheck
			if err != nil {
				if errors.Is(err, codeforces.ErrRatingChangesUnavailable) {
					// Usually means contest was unrated
					continue
				}
				failCnt++
				fmt.Print("\r\033[K") // Clear progress bar line
				log.Printf("Failed to fetch contest: %v\n", err)
				// Sleep before reprinting bar (doesn't want to work without this)
				go func() {
					time.Sleep(100 * time.Millisecond)
					bar.RenderBlank() //nolint:errcheck
				}()
			}
		}

		outputStr := fmt.Sprintf("Fetched %d/%d contests", len(contestIDs)-failCnt, len(contestIDs))
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
