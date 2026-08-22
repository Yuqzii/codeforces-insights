package recommender

import (
	"context"
	"runtime"
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yuqzii/codeforces-insights/internal/codeforces"
)

type mockProblemRepository struct {
	mock.Mock
}

func (m *mockProblemRepository) GetProblemsFromContests(ctx context.Context, ids []int) (
	map[int][]codeforces.Problem, error) {

	args := m.Called(ids)
	return args.Get(0).(map[int][]codeforces.Problem), args.Error(1)
}

func (m *mockProblemRepository) GetProblemsWithTags(ctx context.Context, tags []string, minRat, maxRat int) (
	[]codeforces.Problem, error) {

	args := m.Called(tags)
	return args.Get(0).([]codeforces.Problem), args.Error(1)
}

func TestGetIdxOfTagConcurrentFirstUse(t *testing.T) {
	previousMaxProcs := runtime.GOMAXPROCS(1)
	t.Cleanup(func() {
		runtime.GOMAXPROCS(previousMaxProcs)
	})

	const workerCount = 128

	r := New(nil)
	indices := make([]int, workerCount)
	start := make(chan struct{})
	var ready sync.WaitGroup
	var done sync.WaitGroup
	ready.Add(workerCount)
	done.Add(workerCount)

	// Queue every lookup behind the write lock so they contend on the same tag's first use.
	r.mu.Lock()
	for i := range workerCount {
		go func() {
			defer done.Done()
			<-start
			ready.Done()
			indices[i] = r.getIdxOfTag("dp")
		}()
	}
	close(start)
	ready.Wait()
	runtime.Gosched()
	r.mu.Unlock()
	done.Wait()

	for i, idx := range indices {
		if idx != 0 {
			t.Fatalf("worker %d got tag index %d, want 0", i, idx)
		}
	}
	assert.Equal(t, map[string]int{"dp": 0}, r.tagToIndex)
}

func TestRecommend(t *testing.T) {
	m := new(mockProblemRepository)
	r := New(m)

	emptySet := make(map[int64]struct{})

	t.Run("Filtering input problems", func(t *testing.T) {
		input := []*codeforces.Problem{{
			Name:  "Very cool Problem",
			Index: "A",
			Tags:  []string{"dp", "math"}},
		}

		expected := &ProbWithScore{
			Problem: &codeforces.Problem{
				Name:  "Cool Problem",
				Index: "B",
				Tags:  []string{"dp"},
			},
		}

		mockCall := m.On("GetProblemsWithTags", []string{"dp", "math"}).Return([]codeforces.Problem{
			*input[0], *expected.Problem,
		}, nil)
		defer mockCall.Unset()

		res, err := r.Recommend(context.Background(), input, emptySet, 1, 0, 0)
		assert.Nil(t, err)
		assert.NotZero(t, len(res))
		assert.Equal(t, expected.Problem.Name, res[0].Problem.Name, "Recommended same problem as in input.")
	})

	t.Run("Finding similarity", func(t *testing.T) {
		input := []*codeforces.Problem{
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

		res, err := r.Recommend(context.Background(), input, emptySet, 2, 0, 0)

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

		mockCall := m.On("GetProblemsFromContests", []int{1}).Return(map[int][]codeforces.Problem{
			1: {
				{
					Name:      "Impossible Problem",
					ContestID: 1,
					Index:     "D",
					Tags:      []string{"dp", "brute force"},
				}, {
					Name:      "More Possible Problem",
					ContestID: 1,
					Index:     "B",
					Tags:      []string{"dp", "brute force"},
				}, {
					Name:      "Super Easy Problem",
					ContestID: 1,
					Index:     "A",
					Tags:      []string{"dp", "brute force"},
				}, expected,
			},
		}, nil)
		defer mockCall.Unset()

		p, err := r.findFirstUnsolvedProblem(context.Background(), 1, []string{"B", "A", "D"})

		assert.Nil(t, err)
		assert.Equal(t, *p, expected)
	})

	t.Run("With numbering", func(t *testing.T) {
		expected := codeforces.Problem{
			Name:      "Kniv og Gaffel (Hard version)",
			ContestID: 42,
			Index:     "C2",
			Tags:      []string{"graph", "dsu"},
		}

		mockCall := m.On("GetProblemsFromContests", []int{42}).Return(map[int][]codeforces.Problem{
			42: {
				{
					Name:      "Kniv og Gaffel (Easy version)",
					ContestID: 42,
					Index:     "C1",
					Tags:      []string{"graph", "dsu"},
				}, {
					Name:      "You should solve this",
					ContestID: 42,
					Index:     "A",
					Tags:      []string{"greedy", "brute force"},
				}, {
					Name:      "You should also solve this",
					ContestID: 42,
					Index:     "B",
					Tags:      []string{"greedy", "brute force"},
				},
				expected,
				{
					Name:      "Is this possible? (Easy)",
					ContestID: 42,
					Index:     "D1",
					Tags:      []string{"greedy", "brute force"},
				}, {
					Name:      "Is this possible? (Hard)",
					ContestID: 42,
					Index:     "D2",
					Tags:      []string{"greedy", "brute force"},
				},
			},
		}, nil)
		defer mockCall.Unset()

		p, err := r.findFirstUnsolvedProblem(context.Background(), 42, []string{"B", "A", "D1", "C1"})

		assert.Nil(t, err)
		assert.Equal(t, *p, expected)
	})

	t.Run("Last is unsolved", func(t *testing.T) {
		expected := codeforces.Problem{
			Name:      "Kniv og Gaffel (Hard version)",
			ContestID: 42,
			Index:     "C2",
			Tags:      []string{"graph", "dsu"},
		}

		mockCall := m.On("GetProblemsFromContests", []int{42}).Return(map[int][]codeforces.Problem{
			42: {
				{
					Name:      "Kniv og Gaffel (Easy version)",
					ContestID: 42,
					Index:     "C1",
					Tags:      []string{"graph", "dsu"},
				}, {
					Name:      "You should solve this",
					ContestID: 42,
					Index:     "A",
					Tags:      []string{"greedy", "brute force"},
				}, {
					Name:      "You should also solve this",
					ContestID: 42,
					Index:     "B",
					Tags:      []string{"greedy", "brute force"},
				},
				expected,
			},
		}, nil)
		defer mockCall.Unset()

		p, err := r.findFirstUnsolvedProblem(context.Background(), 42, []string{"A", "B", "C1"})

		assert.Nil(t, err)
		assert.Equal(t, *p, expected)
	})
}

func TestFindSolvedRecentContests(t *testing.T) {
	m := new(mockProblemRepository)
	r := New(m)

	type args struct {
		subs     []codeforces.Submission
		lookback int
	}
	tests := []struct {
		name     string
		args     args
		expected map[int][]*codeforces.Problem
	}{
		{
			name: "Group problems by contests",
			args: args{
				lookback: 1,
				subs: []codeforces.Submission{
					{
						Timestamp: 200,
						ContestID: 100,
						Author:    codeforces.Party{ParticipantType: "CONTESTANT"},
						Problem:   codeforces.Problem{Index: "B"},
					}, {
						Timestamp: 100,
						ContestID: 100,
						Author:    codeforces.Party{ParticipantType: "CONTESTANT"},
						Problem:   codeforces.Problem{Index: "A"},
					},
				},
			},
			expected: map[int][]*codeforces.Problem{
				100: {{Index: "B"}, {Index: "A"}},
			},
		}, {
			name: "Ignore non-contestant submissions",
			args: args{
				lookback: 1,
				subs: []codeforces.Submission{
					{
						Timestamp: 825,
						ContestID: 67,
						Author:    codeforces.Party{ParticipantType: "CONTESTANT"},
						Problem:   codeforces.Problem{Index: "C"},
					}, {
						Timestamp: 1000,
						ContestID: 69,
						Author:    codeforces.Party{ParticipantType: "PRACTICE"},
						Problem:   codeforces.Problem{Index: "D"},
					}, {
						Timestamp: 600,
						ContestID: 67,
						Author:    codeforces.Party{ParticipantType: "CONTESTANT"},
						Problem:   codeforces.Problem{Index: "B"},
					},
				},
			},
			expected: map[int][]*codeforces.Problem{
				67: {{Index: "C"}, {Index: "B"}},
			},
		}, {
			name: "Empty input",
			args: args{
				lookback: 5,
				subs:     []codeforces.Submission{},
			},
			expected: map[int][]*codeforces.Problem{},
		}, {
			name: "Respect lookback",
			args: args{
				lookback: 2,
				subs: []codeforces.Submission{
					{
						Timestamp: 1000,
						ContestID: 420,
						Author:    codeforces.Party{ParticipantType: "CONTESTANT"},
						Problem:   codeforces.Problem{Index: "A"},
					}, {
						Timestamp: 2000,
						ContestID: 421,
						Author:    codeforces.Party{ParticipantType: "CONTESTANT"},
						Problem:   codeforces.Problem{Index: "B"},
					}, {
						Timestamp: 4267,
						ContestID: 670,
						Author:    codeforces.Party{ParticipantType: "CONTESTANT"},
						Problem:   codeforces.Problem{Index: "A"},
					},
				},
			},
			expected: map[int][]*codeforces.Problem{
				421: {{Index: "B"}},
				670: {{Index: "A"}},
			},
		},
	}

	for _, tt := range tests {
		// Update ContestID of argument problems.
		for i := range tt.args.subs {
			tt.args.subs[i].Problem.ContestID = tt.args.subs[i].ContestID
		}

		// Update ContestID of expected problems.
		for i := range tt.expected {
			for j := range tt.expected[i] {
				tt.expected[i][j].ContestID = i
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			actual := r.FindSolvedRecentContests(tt.args.subs, tt.args.lookback)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
