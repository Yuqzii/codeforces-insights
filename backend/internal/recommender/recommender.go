package recommender

import (
	"context"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type ProblemRepo interface {
	GetProblemsFromContest(ctx context.Context, id int) ([]codeforces.Problem, error)
	// Should return all problems matching at least one tag.
	GetProblemsWithTags(ctx context.Context, tag []string) ([]codeforces.Problem, error)
}

type Recommender struct {
	probRepo ProblemRepo
}

func New(repo ProblemRepo) *Recommender {
	return &Recommender{
		probRepo: repo,
	}
}
