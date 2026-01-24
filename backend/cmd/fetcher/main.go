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

	cfTimeBetweenReqs time.Duration = 2100 * time.Millisecond // 2.1 seconds to be nice with CF server

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
		cfTimeBetweenReqs,
	)

	f := fetcher.New(cfClient, db, cfClient, db, db)

	fetchContests := flag.Bool("contests", false, "Should we fetch contests?")
	fetchProblems := flag.Bool("problems", false, "Should we fetch problems?")
	flag.Parse()

	if *fetchContests {
		log.Println("Finding unfetched contests")
		unfetched, err := f.FindUnfetchedContests()
		if err != nil {
			log.Fatalf("Failed to find unfetched contests: %v\n", err)
		}

		log.Printf("Starting fetching for %d contests\n", len(unfetched))
		bar := progressbar.Default(int64(len(unfetched)), "Fetching contests")
		failCnt := 0

		results := fetcher.CreateWorkers(workerCnt, unfetched, cfClient, db, cfClient, db, db)
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
					if err = bar.RenderBlank(); err != nil {
						log.Printf("Failed rendering progress bar: %v", err)
					}
				}()
			}
		}

		outputStr := fmt.Sprintf("Fetched %d/%d contests", len(unfetched)-failCnt, len(unfetched))
		log.Println(outputStr)
	}

	if *fetchProblems {
		log.Println("Fetching problems...")

		err = f.FetchProblems(context.TODO())
		if err != nil {
			log.Printf("Failed to fetch problems: %s", err)
		} else {
			log.Println("Successfully fetched and updated problems")
		}
	}
}
