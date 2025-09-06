package stats

import (
	"slices"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type Tag struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

// Gets the count of each category/tag of the user's solved problems.
// @param solved Slice of Submission to be evaluated. Should use FilterSolved before passing to this.
// @return Sorted slice of Tag based on count.
func SolvedTags(solved []codeforces.Submission) []Tag {
	m := make(map[string]int)
	for _, s := range solved {
		for _, t := range s.Problem.Tags {
			m[t]++
		}
	}

	res := make([]Tag, 0, len(m))
	for t, c := range m {
		res = append(res, Tag{Tag: t, Count: c})
	}

	slices.SortFunc(res, func(a, b Tag) int {
		if a.Count < b.Count {
			return -1
		} else if a.Count > b.Count {
			return 1
		}
		return 0
	})
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
