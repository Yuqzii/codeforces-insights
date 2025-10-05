package codeforces

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Contestant struct {
	Rank          int
	Points        float64
	Penalty       int
	ID            uint64
	OldRating     int
	NewRating     int
	MemberHandles []string
}

type Contest struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"startTime"`
	Duration  int       `json:"durationSeconds"`
	Phase     string    `json:"phase"`
}

var ErrNoStandings = errors.New("could not find standings")

func (c *client) GetContestStandings(ctx context.Context, id int) ([]Contestant, *Contest, error) {
	endpoint := "contest.standings?"
	params := url.Values{}
	params.Set("contestId", strconv.Itoa(id))
	params.Set("from", "1")
	params.Set("showUnofficial", "false")

	resChan := c.makeRequest(ctx, endpoint+params.Encode())
	r := <-resChan
	if r.err != nil {
		return nil, nil, fmt.Errorf("getting contest standings from Codeforces: %w", r.err)
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
	if err := json.Unmarshal(r.body, &apiResp); err != nil {
		return nil, nil, err
	}

	if apiResp.Status != "OK" {
		return nil, nil, fmt.Errorf("%w: %s", ErrCodeforcesReturnedFail, apiResp.Comment)
	}

	return apiResp.Result.Contestants, &apiResp.Result.Contest, nil
}

func (c *client) GetContests(ctx context.Context) ([]Contest, error) {
	endpoint := "contest.list"

	resChan := c.makeRequest(ctx, endpoint)
	r := <-resChan
	if r.err != nil {
		return nil, fmt.Errorf("getting contest list from Codeforces: %w", r.err)
	}

	var apiResp apiResponse[Contest]
	if err := json.Unmarshal(r.body, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("%w: %s", ErrCodeforcesReturnedFail, apiResp.Comment)
	}

	return apiResp.Result, nil
}

func (c *Contestant) UnmarshalJSON(data []byte) error {
	type rawContestant struct {
		Rank    int     `json:"rank"`
		Points  float64 `json:"points"`
		Penalty int     `json:"penalty"`
		Party   struct {
			ParticipantID uint64 `json:"participantId"`
			Members       []struct {
				Handle string `json:"handle"`
			} `json:"members"`
		} `json:"party"`
	}

	var raw rawContestant
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.Rank = raw.Rank
	c.Points = raw.Points
	c.Penalty = raw.Penalty
	c.ID = raw.Party.ParticipantID

	for _, member := range raw.Party.Members {
		c.MemberHandles = append(c.MemberHandles, member.Handle)
	}

	return nil
}

func (c *Contest) UnmarshalJSON(data []byte) error {
	type rawContest struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		StartTime int64  `json:"startTimeSeconds"`
		Duration  int    `json:"durationSeconds"`
		Phase     string `json:"phase"`
	}

	var raw rawContest
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.ID = raw.ID
	c.Name = raw.Name
	c.StartTime = time.Unix(raw.StartTime, 0)
	c.Duration = raw.Duration
	c.Phase = raw.Phase

	return nil
}
