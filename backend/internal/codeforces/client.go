package codeforces

import (
	"context"
	"errors"
	"io"
	"net/http"

	"golang.org/x/time/rate"
)

var ErrCodeforcesReturnedFail = errors.New("Codeforces returned status FAILED") // nolint:staticcheck

type client struct {
	client    *http.Client
	limiter   *rate.Limiter
	url       string
	requests  reqQueue
	receivers map[string][]receiver
}

func NewClient(httpClient *http.Client, url string, requestsPerSecond float64, burst int) *client {
	return &client{
		client:    httpClient,
		limiter:   rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
		url:       url,
		requests:  reqQueue{},
		receivers: make(map[string][]receiver),
	}
}

type receiver struct {
	ctx context.Context
	chn chan<- requestResult
}

type requestResult struct {
	resp *http.Response
	err  error
}

func (c *client) makeRequest(ctx context.Context, method, endpoint string) <-chan requestResult {
	reqs, queued := c.receivers[endpoint]
	if queued {
		// Request is already queued, just add to receivers
		respChan := make(chan requestResult)
		reqs = append(reqs, receiver{
			ctx: ctx,
			chn: respChan,
		})
		return respChan
	}

	// Push request to queue and create receiver list
	c.requests.push(endpoint)
	c.receivers[endpoint] = make([]receiver, 0, 1)

	respChan := make(chan requestResult)
	c.receivers[endpoint] = append(c.receivers[endpoint], receiver{
		ctx: ctx,
		chn: respChan,
	})

	return respChan
}

type apiResponse[T any] struct {
	Status  string `json:"status"`
	Result  []T    `json:"result"`
	Comment string `json:"comment,omitempty"`
}

func closeResponseBody(b io.ReadCloser) {
	_ = b.Close()
}
