package codeforces

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type Submission struct {
	ID                  int     `json:"id"`
	Verdict             string  `json:"verdict"`
	Problem             Problem `json:"problem"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
	Timestamp           int     `json:"creationTimeSeconds"`
}

type Problem struct {
	Name      string   `json:"name" db:"name"`
	ContestID int      `json:"contestId,omitempty" db:"contest_id"`
	Index     string   `json:"index" db:"index"`
	Rating    int      `json:"rating" db:"rating"`
	Tags      []string `json:"tags" db:"tags"`
}

func (c *client) GetSubmissions(ctx context.Context, handle string) ([]Submission, error) {
	endpoint := "user.status?"
	params := url.Values{}
	params.Set("handle", handle)
	params.Set("from", "1")           // Get submissions starting from most recent
	params.Set("count", "1000000000") // Max allowed from Codeforces

	resChan, err := c.makeRequest(ctx, endpoint+params.Encode())
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
		return nil, fmt.Errorf("getting submissions from Codeforces: %w", r.err)
	}

	var apiResp apiResponse[Submission]
	if err := json.Unmarshal(r.body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("%w: %s", ErrCodeforcesReturnedFail, apiResp.Comment)
	}

	return apiResp.Result, nil
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
