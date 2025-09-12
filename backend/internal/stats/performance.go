package stats

import (
	"math"
)

const (
	maxRating     int     = 6000
	minRating     int     = -500
	ratingRange   int     = maxRating - minRating
	defaultRating int     = 1400
	eloScale      float64 = 400
)

func generateEloWinProb() []float64 {
	prob := make([]float64, ratingRange*2+1)
	for i := -ratingRange; i <= ratingRange; i++ {
		prob[i+ratingRange] = 1.0 / (1.0 + math.Pow(10, float64(i)/eloScale)) // Standard ELO calculation
	}
	return prob
}
