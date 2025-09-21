package codeforces

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

type Contestant struct {
	Rank    int     `json:"rank"`
	Points  float64 `json:"points"`
	Penalty int     `json:"penalty"`
	Rating  int
}

type Contest struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	StartTime int    `json:"startTimeSeconds"`
	Duration  int    `json:"durationSeconds"`
}

var ErrNoStandings = errors.New("could not find standings")

func (c *client) GetContestStandings(ctx context.Context, id int) ([]Contestant, *Contest, error) {
	endpoint := "contest.standings?"
	params := url.Values{}
	params.Set("contestId", strconv.Itoa(id))
	params.Set("from", "1")
	params.Set("showUnofficial", "false")

	resp, err := c.makeRequest(ctx, "GET", endpoint+params.Encode())
	if err != nil {
		return nil, nil, fmt.Errorf("getting contest standings from Codeforces: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	// Special apiResponse struct as the Codeforces API returns an unusual json format for this endpoint.
	type apiResponse struct {
		Status string `json:"status"`
		Result struct {
			Contestants []Contestant `json:"rows"`
			Contest     Contest      `json:"contest"`
		} `json:"result"`
		Comment string `json:"comment,omitempty"`
	}
	var apiResp apiResponse
	json.Unmarshal(body, &apiResp)

	if apiResp.Status != "OK" {
		return nil, nil, fmt.Errorf("%w: %s", ErrCodeforcesReturnedFail, apiResp.Comment)
	}

	return apiResp.Result.Contestants, &apiResp.Result.Contest, nil
}
