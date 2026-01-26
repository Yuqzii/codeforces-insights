package recommender

import (
	"math"
	"math/bits"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

const totalTagCnt int = 40 // Actual value is 36, but added some for safety in case new tags are added.

type vec [totalTagCnt]float64

// Because problem vectors are either 1 (has tag) or 0 (does not have tag),
// it can be efficiently represented using a bitset.
type probVec uint64

func (r *recommender) problemToVec(p *codeforces.Problem) probVec {
	var res probVec

	for _, t := range p.Tags {
		idx := r.getIdxOfTag(t)
		res |= 1 << idx
	}

	return res
}

func (r *recommender) tagsToVec(tags []string) vec {
	var res vec

	for _, t := range tags {
		idx := r.getIdxOfTag(t)
		res[idx]++
	}

	return res
}

func similarity(u *vec, v probVec) float64 {
	return dotProduct(u, v) / (u.magnitude() * v.magnitude())
}

func dotProduct(u *vec, v probVec) float64 {
	var res float64

	for i := range totalTagCnt {
		if v&(1<<i) != 0 {
			res += u[i]
		}
	}

	return res
}

func (v *vec) magnitude() float64 {
	var res float64

	for i := range totalTagCnt {
		res += v[i] * v[i]
	}

	return math.Sqrt(res)
}

func (v probVec) magnitude() float64 {
	count := bits.OnesCount64(uint64(v))
	return math.Sqrt(float64(count))
}
