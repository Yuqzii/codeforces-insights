package recommender

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type ProblemRepository interface {
	GetProblemsFromContest(ctx context.Context, id int) ([]codeforces.Problem, error)
	// Should return all problems matching at least one tag.
	GetProblemsWithTags(ctx context.Context, tags []string) ([]codeforces.Problem, error)
}

type recommender struct {
	probRepo ProblemRepository

	tagToIndex map[string]int
	mu         sync.RWMutex
}

func New(repo ProblemRepository) *recommender {
	return &recommender{
		probRepo:   repo,
		tagToIndex: make(map[string]int),
	}
}

var ErrNoUnsolvedProblem = errors.New("there are no unsolved problems for this contest")

// Converts all the problems tags into a vector, and compares the direction of this vector
// to each vector of other problems with at least one tag in common to find the most similar problems.
// @param probs Slice containing problems that we want to find similar problems to.
// @param cnt Amount of problems to recommend.
// @return A slice of length cnt, the recommended problems.
func (r *recommender) Recommend(ctx context.Context, probs []codeforces.Problem, cnt int) ([]*probWithScore, error) {
	tags := make([]string, 0)
	unavailableProbs := make(map[int64]struct{})
	for _, p := range probs {
		tags = append(tags, p.Tags...)
		unavailableProbs[p.Hash()] = struct{}{}
	}

	u := r.tagsToVec(tags)

	// Remove duplicate tags
	slices.Sort(tags)
	tags = slices.Compact(tags)

	probs, err := r.probRepo.GetProblemsWithTags(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("getting problems with tags '%s': %w", tags, err)
	}

	pq := make(probPQ, 0, cnt)
	heap.Init(&pq)

	for _, prob := range probs {
		_, isUnavailable := unavailableProbs[prob.Hash()]
		if isUnavailable {
			continue
		}

		v := r.problemToVec(&prob)
		score := similarity(&u, v)

		ps := probWithScore{
			Score:   score,
			Problem: &prob,
		}

		heap.Push(&pq, &ps)
		// Make sure we only keep cnt best.
		if pq.Len() > cnt {
			heap.Pop(&pq)
		}
	}

	return pq, nil
}

// @param indices Slice of the indices of the solved problems for the contest.
func (r *recommender) FindFirstUnsolvedProblem(ctx context.Context, contestID int,
	indices []string) (*codeforces.Problem, error) {

	allProbs, err := r.probRepo.GetProblemsFromContest(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("getting problems to contest %d: %w", contestID, err)
	}

	// Sort problems by index (A, B, C, D1, D2, ...).
	slices.SortFunc(allProbs, func(a, b codeforces.Problem) int {
		return strings.Compare(a.Index, b.Index)
	})
	slices.Sort(indices)

	for i := range indices {
		if indices[i] != allProbs[i].Index {
			return &allProbs[i], nil
		}
	}

	return nil, ErrNoUnsolvedProblem
}

func (r *recommender) getIdxOfTag(tag string) int {
	r.mu.RLock()
	idx, ok := r.tagToIndex[tag]
	r.mu.RUnlock()
	if ok {
		return idx
	}

	r.mu.Lock()
	r.tagToIndex[tag] = len(r.tagToIndex)
	r.mu.Unlock()

	return r.tagToIndex[tag]
}
