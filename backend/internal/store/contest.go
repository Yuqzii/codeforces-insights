package store

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/db"
)

// First tries getting contest results from the DB, if unsuccessful tries the API.
func (s *Store) GetContestResults(ctx context.Context, id int) (
	[]codeforces.Contestant, *codeforces.Contest, error) {

	contestants, contest, err := s.db.GetContestResults(ctx, id)
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

	ratingMap := make(map[string]*codeforces.RatingChange)
	for i := range ratings {
		ratingMap[ratings[i].Handle] = &ratings[i]
	}

	// Set ratings of contestants
	for i, c := range contestants {
		for _, handle := range c.MemberHandles {
			r, ok := ratingMap[handle]
			// Use rating of party member with maximum previous rating
			if ok && r.OldRating > contestants[i].OldRating {
				contestants[i].OldRating = r.OldRating
				contestants[i].NewRating = r.NewRating
			}
		}
	}

	return contestants, contest, nil
}
