package stats

import (
	"github.com/yuqzii/cf-stats/internal/codeforces"
)

// Gets the count of each category/tag of the user's solved problems
// @param solved Slice of Submission to be evaluated. Should use FilterSolved before passing to this.
func SolvedTags(solved []codeforces.Submission) map[string]int {
	res := make(map[string]int)
	for _, s := range solved {
		for _, t := range s.Problem.Tags {
			res[t]++
		}
	}
	return res
}

// Gets the count of each rating of the user's solved problems
// @param solved Slice of Submission to be evaluated. Should use FilterSolved before passing to this.
func SolvedRatings(solved []codeforces.Submission) map[int]int {
	res := make(map[int]int)
	for _, s := range solved {
		if s.Problem.Rating != 0 {
			res[s.Problem.Rating]++
		}
	}
	return res
}

func FilterSolved(sub []codeforces.Submission) []codeforces.Submission {
	res := make([]codeforces.Submission, 0)
	for _, s := range sub {
		if s.Verdict == "OK" {
			res = append(res, s)
		}
	}
	return res
}
