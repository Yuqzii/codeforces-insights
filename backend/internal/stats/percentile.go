package stats

import (
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type PercentileCalc struct {
	prefix []int
}

func NewPercentile(users []codeforces.User) *PercentileCalc {
	p := &PercentileCalc{
		prefix: make([]int, maxRating+1),
	}

	for _, user := range users {
		p.prefix[user.Rating]++
	}
	for i := len(p.prefix) - 2; i >= 0; i-- {
		p.prefix[i] += p.prefix[i+1]
	}

	return p
}

func (p *PercentileCalc) GetPercentile(rating int) float64 {
	if rating > maxRating || rating < 0 {
		return 0
	}
	return float64(p.prefix[rating]) / float64(p.prefix[maxRating])
}
