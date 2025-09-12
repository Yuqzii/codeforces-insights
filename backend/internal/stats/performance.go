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
	defaultRating int     = 1400
	ratingOffset  int     = -minRating
	eloScale      float64 = 400
)

var eloWinProb = generateEloWinProb()

func generateEloWinProb() []float64 {
	prob := make([]float64, ratingRange*2+1)
	for i := -ratingRange; i <= ratingRange; i++ {
		prob[i+ratingRange] = 1.0 / (1.0 + math.Pow(10, float64(i)/eloScale)) // Standard ELO calculation
	}
	return prob
}

// Computes expected ranks for all possible ratings
func calculateSeed(contestants []codeforces.Contestant) []float64 {
	counts := make([]float64, ratingRange)
	for _, c := range contestants {
		counts[c.Rating+ratingOffset] += 1
	}

	seed := fft.ComplexToFloat(fft.Convolve(fft.FloatToComplex(eloWinProb), fft.FloatToComplex(counts)))

	// Seed base case is 1
	for i := range seed {
		seed[i] += 1
	}

	return seed
}

// Returns expected rank for rating r, excluding one contestant with rating exclude.
func getSeed(seed []float64, r, exclude int) float64 {
	return seed[r+ratingRange+ratingOffset] - eloWinProb[r-exclude+ratingRange]
}
