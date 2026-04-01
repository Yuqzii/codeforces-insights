package stats

import (
	_ "embed"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuqzii/codeforces-insights/internal/codeforces"
	"github.com/yuqzii/codeforces-insights/internal/store"
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

//go:embed testdata/performance_snapshot.json
var snapshotResult []byte

type snapshotEntry struct {
	Rank        int `json:"rank"`
	Rating      int `json:"rating"`
	Performance int `json:"performance"`
}

func TestPerfmanceCalculation(t *testing.T) {
	err := json.Unmarshal(testdataStandings, &testStandings)
	require.Nil(t, err)
	err = json.Unmarshal(testdataRatings, &testRatings)
	require.Nil(t, err)

	store.MapRatingToContestants(testRatings.Ratings, testStandings.Result.Contestants)
	seed := CalculateSeed(testStandings.Result.Contestants, &testStandings.Result.Contest)

	if os.Getenv("UPDATE_SNAPSHOT") == "true" {
		t.Log("Updating performance snapshot")

		entries := make([]snapshotEntry, len(testStandings.Result.Contestants))

		for i, c := range testStandings.Result.Contestants {
			entries[i].Rank = c.Rank
			entries[i].Rating = c.OldRating

			perf := seed.CalculatePerformance(c.Rank, c.OldRating)
			entries[i].Performance = perf
		}

		data, err := json.Marshal(entries)
		require.Nil(t, err, "Failed marshalling snapshot data")

		err = os.WriteFile("testdata/performance_snapshot.json", data, 0644)
		require.Nil(t, err, "Failed writing perfomance snapshot: %v", err)

		t.Log("Successfully updated performance snapshot")
		return
	}

	var snapshot []snapshotEntry
	err = json.Unmarshal(snapshotResult, &snapshot)
	require.Nil(t, err, "Failed parsing snapshot")

	for _, s := range snapshot {
		actual := seed.CalculatePerformance(s.Rank, s.Rating)
		assert.Equal(t, s.Performance, actual)
	}
}
