package codeforces

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	Name      string   `json:"name"`
	ContestID int      `json:"contestId,omitempty"`
	Index     string   `json:"index"`
	Rating    int      `json:"rating"`
	Tags      []string `json:"tags"`
}

func (c *client) GetSubmissions(ctx context.Context, handle string) ([]Submission, error) {
	endpoint := "user.status?"
	params := url.Values{}
	params.Set("handle", handle)
	params.Set("from", "1")           // Get submissions starting from most recent
	params.Set("count", "1000000000") // Max allowed from Codeforces

	resp, err := c.makeRequest(ctx, "GET", endpoint+params.Encode())
	if err != nil {
		return nil, fmt.Errorf("getting submissions from Codeforces: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp apiResponse[Submission]
	json.Unmarshal(body, &apiResp)

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("%w: %s", ErrCodeforcesReturnedFail, apiResp.Comment)
	}

	return apiResp.Result, nil
}
