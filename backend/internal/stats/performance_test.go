package stats

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/store"
)

//go:embed testdata/contest_standings.json
var testdataStandings []byte

//go:embed testdata/contest_ratings.json
var testdataRatings []byte

var testStandings struct {
	Result struct {
		Contest     codeforces.Contest      `json:"contest"`
		Contestants []codeforces.Contestant `json:"rows"`
	} `json:"result"`
}

var testRatings struct {
	Ratings []codeforces.RatingChange `json:"result"`
}

func BenchmarkPerformanceCalculation(b *testing.B) {
	err := json.Unmarshal(testdataStandings, &testStandings)
	require.Nil(b, err)
	err = json.Unmarshal(testdataRatings, &testRatings)
	require.Nil(b, err)

	store.MapRatingToContestants(testRatings.Ratings, testStandings.Result.Contestants)

	b.ResetTimer()

	for b.Loop() {
		seed := CalculateSeed(testStandings.Result.Contestants, &testStandings.Result.Contest)
		seed.CalculatePerformance(4000, 1200)
	}
}
