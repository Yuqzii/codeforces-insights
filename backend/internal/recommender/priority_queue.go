package recommender

import "github.com/yuqzii/cf-stats/internal/codeforces"

type ProbWithScore struct {
	Score   float64             `json:"score"`
	Problem *codeforces.Problem `json:"problem"`
	index   int
}

type probPQ []*ProbWithScore

func (pq probPQ) Len() int {
	return len(pq)
}

func (pq probPQ) Less(i, j int) bool {
	return pq[i].Score < pq[j].Score
}

func (pq probPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *probPQ) Push(x any) {
	p := x.(*ProbWithScore)
	p.index = len(*pq)
	*pq = append(*pq, p)
}

func (pq *probPQ) Pop() any {
	last := (*pq)[len(*pq)-1]
	(*pq)[len(*pq)-1] = nil
	*pq = (*pq)[0 : len(*pq)-1]
	return last
}
