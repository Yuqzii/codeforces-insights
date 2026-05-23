package recommender

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/yuqzii/codeforces-insights/internal/codeforces"
	"github.com/yuqzii/codeforces-insights/internal/db"
)

type ProblemRepository interface {
	GetProblemsFromContest(ctx context.Context, id int) ([]codeforces.Problem, error)
	// Should return all problems matching at least one tag.
	GetProblemsWithTags(ctx context.Context, tags []string, minRat, maxRat int) (
		[]codeforces.Problem, error)
}

type recommender struct {
	probRepo ProblemRepository

	tagToIndex map[string]int
	mu         sync.RWMutex
}

func New(repo ProblemRepository) *recommender {
	return &recommender{
		probRepo:   repo,
		tagToIndex: make(map[string]int),
	}
}

var (
	ErrNoUnsolvedProblem = errors.New("there are no unsolved problems for this contest")
	ErrInvalidIndices    = errors.New("given problem indices are invalid")
)

// Converts all the problems tags into a vector, and compares the direction of this vector
// to each vector of other problems with at least one tag in common to find the most similar problems.
// @param probs Slice containing problems that we want to find similar problems to.
// @param disallowedProbs Map of empty structs (effectively a set) containing hashes of problems
// to not recommend. E.g. problems the user has already solved.
// @param cnt Amount of problems to recommend.
// @param minRat Minimum rating of recommended problems.
// @param maxRat Maximum rating of recommended problems.
// @return A slice of length cnt, the recommended problems.
func (r *recommender) Recommend(ctx context.Context, probs []*codeforces.Problem,
	disallowedProbs map[int64]struct{}, cnt, minRat, maxRat int) ([]*ProbWithScore, error) {

	tags := make([]string, 0)
	for _, p := range probs {
		tags = append(tags, p.Tags...)
		disallowedProbs[p.Hash()] = struct{}{}
	}

	u := r.tagsToVec(tags)

	// Remove duplicate tags
	slices.Sort(tags)
	tags = slices.Compact(tags)

	allProbs, err := r.probRepo.GetProblemsWithTags(ctx, tags, minRat, maxRat)
	if err != nil {
		return nil, fmt.Errorf("getting problems with tags '%s': %w", tags, err)
	}

	pq := make(probPQ, 0, cnt)
	heap.Init(&pq)

	for _, prob := range allProbs {
		_, isDisallowed := disallowedProbs[prob.Hash()]
		if isDisallowed {
			continue
		}

		v := r.problemToVec(&prob)
		score := similarity(&u, v) * sigmoid(float64(prob.ContestID))

		ps := ProbWithScore{
			Score:   score,
			Problem: &prob,
		}

		heap.Push(&pq, &ps)
		// Make sure we only keep cnt best.
		if pq.Len() > cnt {
			heap.Pop(&pq)
		}
	}

	return pq, nil
}

// @param solvedByContest Map with key as contest ID and value as slice of problems.
// This should generally be the output of FindSolvedRecentContests.
func (r *recommender) FindUnsolvedProblems(ctx context.Context,
	solvedByContest map[int][]*codeforces.Problem) ([]*codeforces.Problem, error) {

	unsolved := make([]*codeforces.Problem, 0, len(solvedByContest))

	for contestID, probs := range solvedByContest {
		indices := make([]string, 0)
		for _, p := range probs {
			indices = append(indices, p.Index)
		}

		unsolvedProb, err := r.findFirstUnsolvedProblem(ctx, contestID, indices)
		if err != nil {
			if errors.Is(err, ErrNoUnsolvedProblem) {
				continue
			}

			if errors.Is(err, db.ErrNoProblemsForContest) {
				log.Printf("Couldn't find problems from contest %d: %v\n", contestID, err)
				continue
			}

			if errors.Is(err, ErrInvalidIndices) {
				return nil, fmt.Errorf("finding unsolved problems: %w", err)
			}

			log.Printf("Error finding unsolved problem for contest %d and indices %v: %v\n",
				contestID, indices, err)
		}

		unsolved = append(unsolved, unsolvedProb)
	}

	return unsolved, nil
}

func (r *recommender) FindSolvedRecentContests(subs []codeforces.Submission,
	lookback int) map[int][]*codeforces.Problem {

	// Sort submissions in descending order of submission time.
	slices.SortFunc(subs, func(a, b codeforces.Submission) int {
		return b.Timestamp - a.Timestamp
	})

	probsByContest := make(map[int][]*codeforces.Problem, lookback)
	for i := range subs {
		if subs[i].Author.ParticipantType != "CONTESTANT" {
			// Submission was not in contest.
			continue
		}

		s, contestStored := probsByContest[subs[i].ContestID]
		if contestStored {
			s = append(s, &subs[i].Problem)
			probsByContest[subs[i].ContestID] = s
		} else {
			if len(probsByContest) >= lookback {
				break
			} else {
				s = make([]*codeforces.Problem, 1)
				s[0] = &subs[i].Problem
				probsByContest[subs[i].ContestID] = s
			}
		}
	}

	return probsByContest
}

func (r *recommender) GetSolvedProblemHashes(ctx context.Context, probs []*codeforces.Problem) ([]int64, error) {
	probsByContest := make(map[int][]codeforces.Problem)
	hashes := make([]int64, len(probs))

	for i, p := range probs {
		contestProbs, ok := probsByContest[p.ContestID]
		if !ok {
			var err error
			contestProbs, err = r.probRepo.GetProblemsFromContest(ctx, p.ContestID)
			if err != nil {
				return nil, fmt.Errorf("getting problems from contest %d: %w", p.ContestID, err)
			}

			probsByContest[p.ContestID] = contestProbs
		}

		j := sort.Search(len(contestProbs), func(i int) bool {
			return contestProbs[i].Index >= p.Index
		})

		actualProb := contestProbs[j]
		hashes[i] = actualProb.Hash()
	}

	return hashes, nil
}

// _𜰾𜰱‾ used as a scalar by contestID such that newer problems gets recommended more
func sigmoid(x float64) float64 {
	const (
		inflection float64 = 1500
		growth     float64 = 500
	)
	return 1.0 / (1.0 + math.Exp(-(x-inflection)/growth))
}

// @param indices Slice of the indices of the solved problems for the contest.
func (r *recommender) findFirstUnsolvedProblem(ctx context.Context, contestID int,
	indices []string) (*codeforces.Problem, error) {

	allProbs, err := r.probRepo.GetProblemsFromContest(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("getting problems to contest %d: %w", contestID, err)
	}

	if len(indices) > len(allProbs) {
		return nil, ErrInvalidIndices
	}

	// Sort problems by index (A, B, C, D1, D2, ...).
	slices.SortFunc(allProbs, func(a, b codeforces.Problem) int {
		return strings.Compare(a.Index, b.Index)
	})
	slices.Sort(indices)

	for i := range indices {
		if indices[i] != allProbs[i].Index {
			return &allProbs[i], nil
		}
	}

	if len(indices) == len(allProbs) {
		return nil, ErrNoUnsolvedProblem
	}

	return &allProbs[len(indices)], nil
}

func (r *recommender) getIdxOfTag(tag string) int {
	r.mu.RLock()
	idx, ok := r.tagToIndex[tag]
	r.mu.RUnlock()
	if ok {
		return idx
	}

	r.mu.Lock()
	r.tagToIndex[tag] = len(r.tagToIndex)
	r.mu.Unlock()

	return r.tagToIndex[tag]
}
