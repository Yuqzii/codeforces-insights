package recommender

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type mockProblemRepository struct {
	mock.Mock
}

func (m *mockProblemRepository) GetProblemsFromContest(ctx context.Context, id int) (
	[]codeforces.Problem, error) {

	args := m.Called(id)
	return args.Get(0).([]codeforces.Problem), args.Error(1)
}

func (m *mockProblemRepository) GetProblemsWithTags(ctx context.Context, tags []string) (
	[]codeforces.Problem, error) {

	args := m.Called(tags)
	return args.Get(0).([]codeforces.Problem), args.Error(1)
}

func TestFindFirstUnsolvedProblem(t *testing.T) {
	m := new(mockProblemRepository)
	r := New(m)

	t.Run("No numbering", func(t *testing.T) {
		expected := codeforces.Problem{
			Name:      "Maybe Impossible Problem",
			ContestID: 1,
			Index:     "C",
			Rating:    2000,
			Tags:      []string{"dp", "brute force"},
		}

		m.On("GetProblemsFromContest", 1).Return([]codeforces.Problem{
			{
				Name:      "Impossible Problem",
				ContestID: 1,
				Index:     "D",
				Rating:    3000,
				Tags:      []string{"dp", "brute force"},
			}, {
				Name:      "More Possible Problem",
				ContestID: 1,
				Index:     "B",
				Rating:    1200,
				Tags:      []string{"dp", "brute force"},
			}, {
				Name:      "Super Easy Problem",
				ContestID: 1,
				Index:     "A",
				Rating:    800,
				Tags:      []string{"dp", "brute force"},
			}, expected,
		}, nil)

		p, err := r.FindFirstUnsolvedProblem(context.Background(), 1, []string{"B", "A", "D"})
		assert.Nil(t, err)
		assert.Equal(t, *p, expected)
	})
}
