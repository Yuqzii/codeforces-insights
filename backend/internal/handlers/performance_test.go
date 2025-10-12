package handlers

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/store"
)

//go:embed testdata/contest_standings.json
var testdataStandings []byte

//go:embed testdata/contest_ratings.json
var testdataRatings []byte

type mockCRP struct {
	mock.Mock
}

var testStandings struct {
	Result struct {
		Contest     codeforces.Contest      `json:"contest"`
		Contestants []codeforces.Contestant `json:"rows"`
	} `json:"result"`
}

var testRatings struct {
	Ratings []codeforces.RatingChange `json:"result"`
}

func (m *mockCRP) GetContestResults(ctx context.Context, id int) (
	[]codeforces.Contestant, *codeforces.Contest, error) {

	return testStandings.Result.Contestants, &testStandings.Result.Contest, nil
}

func BenchmarkPerformanceCalculation(b *testing.B) {
	err := json.Unmarshal(testdataStandings, &testStandings)
	require.Nil(b, err)
	err = json.Unmarshal(testdataRatings, &testRatings)
	require.Nil(b, err)

	store.MapRatingToContestants(testRatings.Ratings, testStandings.Result.Contestants)

	mock := new(mockCRP)
	p := perfManager{
		jobs: make(chan perfJob, 10),
		crp:  mock,
	}

	go p.perfWorker()

	r := codeforces.RatingChange{
		ContestID: testStandings.Result.Contest.ID,
		Rank:      4000,
		OldRating: 1200,
		NewRating: 1300,
		Handle:    "testUser",
	}

	b.ResetTimer()

	for b.Loop() {
		resChn := make(chan perfResult)
		p.addJob(context.Background(), &r, resChn)

		res := <-resChn
		assert.Nil(b, res.err)
	}
}
