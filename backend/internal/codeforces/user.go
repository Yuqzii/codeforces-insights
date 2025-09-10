package codeforces

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
)

type User struct {
	Handle     string `json:"handle"`
	Rating     int    `json:"rating"`
	MaxRating  int    `json:"maxRating"`
	Rank       string `json:"rank"`
	MaxRank    string `json:"maxRank"`
	TitlePhoto string `json:"titlePhoto"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Country    string `json:"country,omitempty"`
}

var ErrUserNotFound = errors.New("could not find a user with that handle")

func (c *client) GetUser(ctx context.Context, handle string) (*User, error) {
	endpoint := "user.info?"
	params := url.Values{}
	params.Set("handles", handle)

	resp, err := c.makeRequest(ctx, "GET", endpoint+params.Encode())
	if err != nil {
		return nil, fmt.Errorf("getting user from Codeforces: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp apiResponse[User]
	json.Unmarshal(body, &apiResp)

	if apiResp.Status != "OK" {
		return nil, fmt.Errorf("%w: %s", ErrCodeforcesReturnedFail, apiResp.Comment)
	}

	if len(apiResp.Result) == 0 {
		return nil, ErrUserNotFound
	}

	return &apiResp.Result[0], nil
}
