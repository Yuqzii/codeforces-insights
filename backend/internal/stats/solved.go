package stats

import (
	"context"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

// Gets the count of each category/tag of the user's solved problems
func (s *service) Categories(handle string) (map[string]int, error) {
	sub, err := s.client.GetSubmissions(context.TODO(), handle)
	if err != nil {
		return nil, err
	}

	solved := filterSolved(sub)
	res := make(map[string]int)
	for _, s := range solved {
		for _, t := range s.Problem.Tags {
			res[t]++
		}
	}
	return res, nil
}

// Gets the count of each rating of the user's solved problems
func (s *service) Ratings(handle string) (map[int]int, error) {
	sub, err := s.client.GetSubmissions(context.TODO(), handle)
	if err != nil {
		return nil, err
	}

	solved := filterSolved(sub)
	res := make(map[int]int)
	for _, s := range solved {
		res[s.Problem.Rating]++
	}
	return res, nil
}

func filterSolved(sub []codeforces.Submission) []codeforces.Submission {
	res := make([]codeforces.Submission, 0)
	for _, s := range sub {
		if s.Verdict == "OK" {
			res = append(res, s)
		}
	}
	return res
}
