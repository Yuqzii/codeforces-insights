package fetcher

import (
	"context"
	"fmt"
)

func (s *Service) FetchProblems(ctx context.Context) error {
	probs, err := s.problemProvider.GetProblems(ctx)
	if err != nil {
		return fmt.Errorf("getting problems: %w", err)
	}

	n := len(probs)
	for i := 0; i < n; i++ {
		isIncomplete := len(probs[i].Tags) == 0 || probs[i].Rating == 0
		if isIncomplete {
			// Swap problem with the last one and reevaluate for efficient removal at the end.
			probs[i], probs[n-1] = probs[n-1], probs[i]
			i--
			n--
		}
	}

	probs = probs[:n] // Remove all incomplete problems.

	err = s.problemRepo.UpsertProblemsBatch(ctx, probs)
	if err != nil {
		return fmt.Errorf("upserting %d problems: %w", n, err)
	}

	return nil
}
