package codeforces

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/time/rate"
)

var ErrCodeforcesReturnedFail = errors.New("Codeforces returned status FAILED") // nolint:staticcheck

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

func closeResponseBody(b io.ReadCloser) {
	_ = b.Close()
}
