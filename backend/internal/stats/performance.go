package stats

import (
	"math"
	"strings"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/fft"
)

const (
	maxRating    int     = 6000
	minRating    int     = -500
	ratingRange  int     = maxRating - minRating
	ratingOffset int     = -minRating
	eloScale     float64 = 400

	avgRatingWeight float64 = 0.4
	div1Rating      int     = 1300
	div2Rating      int     = 1175
	div3Rating      int     = 1100
	div4Rating      int     = 1050
)

type ContestSeed struct {
	seed []float64
}

var eloWinProb = generateEloWinProb()

func (s *ContestSeed) CalculatePerformance(rank, rating int) int {
	perf := s.rankToRating(float64(rank), rating)
	return perf
}

func generateEloWinProb() []float64 {
	prob := make([]float64, ratingRange*2+1)
	for i := -ratingRange; i <= ratingRange; i++ {
		prob[i+ratingRange] = 1.0 / (1.0 + math.Pow(10, float64(i)/eloScale)) // Standard ELO calculation
	}
	return prob
}

// Computes expected ranks for all possible ratings.
func CalculateSeed(contestants []codeforces.Contestant, contest *codeforces.Contest) *ContestSeed {
	// Set default rating based on contest division
	defaultRating := div2Rating
	if strings.Contains(contest.Name, "Div. 3") {
		defaultRating = div3Rating
	} else if strings.Contains(contest.Name, "Div. 4") {
		defaultRating = div4Rating
	} else if strings.Contains(contest.Name, "Div. 1") && !strings.Contains(contest.Name, "Div. 2") {
		defaultRating = div1Rating
	}

	// Calculate average rating
	ratSum, cnt := 0.0, 0
	for i := range contestants {
		if contestants[i].OldRating != 0 {
			ratSum += float64(contestants[i].OldRating)
			cnt++
		}
	}
	ratAvg := ratSum / float64(cnt)

	// Calculate default rating combined with average
	defaultRating = int((float64(defaultRating) + ratAvg*avgRatingWeight) / (1.0 + avgRatingWeight))

	counts := make([]float64, ratingRange)
	for _, c := range contestants {
		if c.OldRating == 0 {
			c.OldRating = defaultRating
		}
		counts[c.OldRating+ratingOffset] += 1
	}

	seedComplex := fft.Convolve(fft.FloatToComplex(eloWinProb), fft.FloatToComplex(counts))
	seed := ContestSeed{
		seed: fft.ComplexToFloat(seedComplex),
	}

	// Seed base case is 1
	for i := range seed.seed {
		seed.seed[i] += 1
	}

	return &seed
}

// Uses binary search to find what rating gives rank targetRank.
// @param rating Rating to exclude from expected rank calculation, rating of contestant we are calculating for.
func (s *ContestSeed) rankToRating(targetRank float64, rating int) int {
	l, r := 2, maxRating
	for l < r {
		mid := (l + r) / 2
		expectedRank := s.get(mid, rating)
		if expectedRank > targetRank {
			// Rating mid gives too high rank
			l = mid + 1
		} else {
			// Rating mid gives too low rank
			r = mid
		}
	}
	// l is now first rating where expected rank < targetRank.
	// - 1 to find last rating where expected rank >= targetRank.
	return l - 1
}

// Returns expected rank for rating r, excluding one contestant with rating exclude.
func (s *ContestSeed) get(r, exclude int) float64 {
	return s.seed[r+ratingRange+ratingOffset] - eloWinProb[r-exclude+ratingRange]
}
