package stats

import (
	"math"

	"github.com/yuqzii/cf-stats/internal/codeforces"
	"github.com/yuqzii/cf-stats/internal/fft"
)

const (
	maxRating     int     = 6000
	minRating     int     = -500
	ratingRange   int     = maxRating - minRating
	ratingOffset  int     = -minRating
	defaultRating int     = 1400
	eloScale      float64 = 400
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
func CalculateSeed(contestants []codeforces.Contestant) *ContestSeed {
	counts := make([]float64, ratingRange)
	for _, c := range contestants {
		if c.Rating == 0 {
			c.Rating = defaultRating
		}
		counts[c.Rating+ratingOffset] += 1
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
