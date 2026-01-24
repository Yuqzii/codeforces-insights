package recommender

import (
	"context"
	"slices"
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

func TestRecommend(t *testing.T) {
	m := new(mockProblemRepository)
	r := New(m)

	t.Run("Filtering input problems", func(t *testing.T) {
		input := []codeforces.Problem{{
			Name:  "Very cool Problem",
			Index: "A",
			Tags:  []string{"dp", "math"}},
		}

		expected := &probWithScore{
			Problem: &codeforces.Problem{
				Name:  "Cool Problem",
				Index: "B",
				Tags:  []string{"dp"},
			},
		}

		mockCall := m.On("GetProblemsWithTags", []string{"dp", "math"}).Return([]codeforces.Problem{
			input[0], *expected.Problem,
		}, nil)
		defer mockCall.Unset()

		res, err := r.Recommend(context.Background(), input, 1)
		assert.Nil(t, err)
		assert.NotZero(t, len(res))
		assert.Equal(t, expected.Problem.Name, res[0].Problem.Name, "Recommended same problem as in input.")
	})

	t.Run("Finding similarity", func(t *testing.T) {
		input := []codeforces.Problem{
			{Index: "A", Tags: []string{"greedy", "graph"}},
		}

		expected := []codeforces.Problem{
			{Index: "C1", Tags: []string{"greedy", "dp"}},
			{Index: "C2", Tags: []string{"graph", "greedy"}},
		}

		allProbs := []codeforces.Problem{
			{Index: "B", Tags: []string{"math", "dsu"}},
		}
		allProbs = append(allProbs, expected...)
		mockCall := m.On("GetProblemsWithTags", []string{"graph", "greedy"}).Return(allProbs, nil)
		defer mockCall.Unset()

		res, err := r.Recommend(context.Background(), input, 2)

		assert.Nil(t, err)
		assert.Equal(t, len(res), 2, "Did not recommend the correct amount of problems.")

		resProbs := make([]*codeforces.Problem, 0, 2)
		for _, ps := range res {
			resProbs = append(resProbs, ps.Problem)
		}

		slices.SortFunc(resProbs, func(a, b *codeforces.Problem) int {
			return int(a.Hash() - b.Hash())
		})
		slices.SortFunc(expected, func(a, b codeforces.Problem) int {
			return int(a.Hash() - b.Hash())
		})

		for i := range len(expected) {
			assert.Equal(t, expected[i], *resProbs[i])
		}
	})
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

		mockCall := m.On("GetProblemsFromContest", 1).Return([]codeforces.Problem{
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
		defer mockCall.Unset()

		p, err := r.FindFirstUnsolvedProblem(context.Background(), 1, []string{"B", "A", "D"})

		assert.Nil(t, err)
		assert.Equal(t, *p, expected)
	})

	t.Run("With numbering", func(t *testing.T) {
		expected := codeforces.Problem{
			Name:      "Kniv og Gaffel (Hard version)",
			ContestID: 42,
			Index:     "C2",
			Rating:    1800,
			Tags:      []string{"graph", "dsu"},
		}

		mockCall := m.On("GetProblemsFromContest", 42).Return([]codeforces.Problem{
			{
				Name:      "Kniv og Gaffel (Easy version)",
				ContestID: 42,
				Index:     "C1",
				Rating:    1400,
				Tags:      []string{"graph", "dsu"},
			}, {
				Name:      "You should solve this",
				ContestID: 42,
				Index:     "A",
				Rating:    900,
				Tags:      []string{"greedy", "brute force"},
			}, {
				Name:      "You should also solve this",
				ContestID: 42,
				Index:     "B",
				Rating:    1100,
				Tags:      []string{"greedy", "brute force"},
			},
			expected,
			{
				Name:      "Is this possible? (Easy)",
				ContestID: 42,
				Index:     "D1",
				Rating:    2500,
				Tags:      []string{"greedy", "brute force"},
			}, {
				Name:      "Is this possible? (Hard)",
				ContestID: 42,
				Index:     "D2",
				Rating:    2900,
				Tags:      []string{"greedy", "brute force"},
			},
		}, nil)
		defer mockCall.Unset()

		p, err := r.FindFirstUnsolvedProblem(context.Background(), 42, []string{"B", "A", "D1", "C1"})

		assert.Nil(t, err)
		assert.Equal(t, *p, expected)
	})
}
