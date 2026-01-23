package recommender

import (
	"context"
	"sync"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type ProblemRepo interface {
	GetProblemsFromContest(ctx context.Context, id int) ([]codeforces.Problem, error)
	// Should return all problems matching at least one tag.
	GetProblemsWithTags(ctx context.Context, tag []string) ([]codeforces.Problem, error)
}

type recommender struct {
	probRepo ProblemRepo

	tagToIndex map[string]int
	mu         sync.RWMutex
}

func New(repo ProblemRepo) *recommender {
	return &recommender{
		probRepo: repo,
	}
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
