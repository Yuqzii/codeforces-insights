package codeforces

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
)

type RatingChange struct {
	ContestID   int    `json:"contestId"`
	ContestName string `json:"contestName"`
	Rank        int    `json:"rank"`
	Timestamp   int    `json:"ratingUpdateTimeSeconds"`
	OldRating   int    `json:"oldRating"`
	NewRating   int    `json:"newRating"`
}

var ErrNoRating = errors.New("user does not have rating")

func (c *client) GetRatingChanges(ctx context.Context, handle string) ([]RatingChange, error) {
	endpoint := "user.rating?"
	params := url.Values{}
	params.Set("handle", handle)

	resp, err := c.makeRequest(ctx, "GET", endpoint+params.Encode())
	if err != nil {
		return nil, fmt.Errorf("getting rating from Codeforces: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp apiResponse[RatingChange]
	json.Unmarshal(body, &apiResp)

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("%w: %s", ErrCodeforcesReturnedFail, apiResp.Comment)
	}

	if len(apiResp.Result) == 0 {
		return nil, ErrNoRating
	}

	return apiResp.Result, nil
}
