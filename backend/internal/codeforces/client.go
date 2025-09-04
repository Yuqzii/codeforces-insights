package codeforces

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/time/rate"
)

var (
	ErrCodeforcesReturnedFail = errors.New("Codeforces returned status FAILED")
	ErrUserNotFound           = errors.New("could not find a user with that handle")
)

type client struct {
	client  *http.Client
	limiter *rate.Limiter
	url     string
}

func NewClient(httpClient *http.Client, url string, requestsPerSecond float64, burst int) *client {
	return &client{
		client:  httpClient,
		limiter: rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
		url:     url,
	}
}

func (c *client) makeRequest(ctx context.Context, method, endpoint string) (*http.Response, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("waiting for limiter: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.url+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	return c.client.Do(req)
}

type apiResponse[T any] struct {
	Status  string `json:"status"`
	Result  []T    `json:"result"`
	Comment string `json:"comment,omitempty"`
}

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
