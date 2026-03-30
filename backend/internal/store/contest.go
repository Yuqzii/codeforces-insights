package store

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/yuqzii/codeforces-insights/internal/codeforces"
	"github.com/yuqzii/codeforces-insights/internal/db"
)

// First tries getting contest results from the DB, if unsuccessful tries the API.
func (s *Store) GetContestResults(ctx context.Context, id int) (
	[]codeforces.Contestant, *codeforces.Contest, error) {

	contestants, contest, err := s.db.GetContestResults(ctx, id, false)
	if err == nil {
		return contestants, contest, nil
	}

	if errors.Is(err, context.Canceled) {
		return nil, nil, err
	}

	if !errors.Is(err, db.ErrContestNotStored) {
		log.Printf("unexpected error querying db: %v\ntrying api", err)
	}

	contestants, contest, err = s.api.GetContestStandings(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("getting contest standings from api: %w", err)
	}

	ratings, err := s.api.GetContestRatingChanges(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("getting contest ratings from api: %w", err)
	}

	MapRatingToContestants(ratings, contestants)

	return contestants, contest, nil
}

// First tries the DB, if unsuccessful tries the API.
// NOTE: can return partially incomplete Contestant data,
// as the Codeforces API does not return all the information we have stored.
func (s *Store) GetContestResultsFromHandle(ctx context.Context, handle string) (
	[]codeforces.Contestant, error) {

	contestants, err := s.db.GetContestResultsFromHandle(ctx, handle)
	if err == nil {
		return contestants, err
	}

	if errors.Is(err, context.Canceled) {
		return nil, err
	}

	if !errors.Is(err, db.ErrNoResultsWithHandle) {
		log.Printf("unexpected error querying db for contest results from handle: %v\ntrying api", err)
	}

	ratings, err := s.api.GetRatingChanges(ctx, handle)
	if err != nil {
		return nil, fmt.Errorf("getting rating changes for handle '%s' from api:, %w", handle, err)
	}

	contestants = make([]codeforces.Contestant, 0, len(ratings))
	for _, r := range ratings {
		contestants = append(contestants, codeforces.Contestant{
			Rank:          r.Rank,
			OldRating:     r.OldRating,
			NewRating:     r.NewRating,
			MemberHandles: []string{r.Handle},
		})
	}

	return contestants, nil
}

// Adds the correct rating to each contestant. (Changes the provided contestants).
func MapRatingToContestants(ratings []codeforces.RatingChange, contestants []codeforces.Contestant) {
	ratingMap := make(map[string]*codeforces.RatingChange)
	for i := range ratings {
		ratingMap[ratings[i].Handle] = &ratings[i]
	}

	// Set ratings of contestants
	for i := range contestants {
		for _, handle := range contestants[i].MemberHandles {
			r, ok := ratingMap[handle]
			// Use rating of party member with maximum previous rating
			if ok && r.OldRating > contestants[i].OldRating {
				contestants[i].OldRating = r.OldRating
				contestants[i].NewRating = r.NewRating
			}
		}
	}
}
