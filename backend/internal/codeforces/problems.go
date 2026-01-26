package codeforces

import (
	"context"
	"encoding/json"
	"fmt"
)

type Problem struct {
	Name      string   `json:"name" db:"name"`
	ContestID int      `json:"contestId,omitempty" db:"contest_id"`
	Index     string   `json:"index" db:"index"`
	Rating    int      `json:"rating" db:"rating"`
	Tags      []string `json:"tags" db:"tags"`
}

func (c *client) GetProblems(ctx context.Context) ([]Problem, error) {
	endpoint := "problemset.problems"

	resChan, err := c.makeRequest(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	var r requestResult
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r = <-resChan:
	}

	if r.err != nil {
		return nil, fmt.Errorf("getting problems from Codeforces: %w", err)
	}

	// Special apiResponse struct as the Codeforces API returns an unusual json format for this endpoint.
	type apiResponse struct {
		Status string `json:"status"`
		Result struct {
			Problems []Problem `json:"problems"`
		} `json:"result"`
		Comment string `json:"comment,omitempty"`
	}
	var apiResp apiResponse
	if err := json.Unmarshal(r.body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("%w: %s", ErrCodeforcesReturnedFail, apiResp.Comment)
	}

	return apiResp.Result.Problems, nil
}

func (p *Problem) Hash() int64 {
	return HashProblemIndex(p.ContestID, p.Index)
}

func HashProblemIndex(contestID int, index string) int64 {
	var res int64
	res = int64(contestID) << 32
	res |= 1 << (index[0] - 'A' + 2) // One bit for indicating the character

	if len(index) > 1 {
		res |= int64(index[1] - '1') // Last two bits indicating number
	}

	return res
}
