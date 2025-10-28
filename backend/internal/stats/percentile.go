package stats

import (
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type percentileCalc struct {
	prefix []int
}

func NewPercentile(users []codeforces.User) *percentileCalc {
	p := &percentileCalc{
		prefix: make([]int, maxRating),
	}

	for _, user := range users {
		p.prefix[user.Rating]++
	}
	for i := len(p.prefix) - 2; i >= 0; i-- {
		p.prefix[i] += p.prefix[i+1]
	}

	return p
}
