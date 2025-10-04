package codeforces

import (
	"context"
	"encoding"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrCodeforcesReturnedFail = errors.New("Codeforces returned status FAILED") // nolint:staticcheck
	ErrAllReceiversCancelled  = errors.New("all receivers to request cancelled")
)

type client struct {
	client          *http.Client
	url             string
	timeBetweenReqs time.Duration

	requests  reqQueue
	receivers map[string][]receiver
}

func NewClient(httpClient *http.Client, url string, timeBetweenReqs time.Duration) *client {
	return &client{
		client:          httpClient,
		timeBetweenReqs: timeBetweenReqs,
		url:             url,
		requests:        reqQueue{},
		receivers:       make(map[string][]receiver),
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

type apiResponse[T any] struct {
	Status  string `json:"status"`
	Result  []T    `json:"result"`
	Comment string `json:"comment,omitempty"`
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

func (c *client) sendNextRequest() error {
	endpoint, err := c.requests.front()
	if err != nil {
		return fmt.Errorf("getting next request: %w", err)
	}
	if err = c.requests.pop(); err != nil {
		return fmt.Errorf("popping request queue: %w", err)
	}

	if c.receiversCancelled(endpoint) {
		return ErrAllReceiversCancelled
	}

	req, err := http.NewRequest("GET", c.url+endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("requesting '%s' from Codeforces: %w", endpoint, err)
	}

	result := requestResult{
		resp: resp,
		err:  nil,
	}
	for _, recvr := range c.receivers[endpoint] {
		recvr.chn <- result
		close(recvr.chn)
	}

	return nil
}

// Returns true if all receivers to endpoint has cancelled their context.
func (c *client) receiversCancelled(endpoint string) bool {
	for _, r := range c.receivers[endpoint] {
		select {
		case <-r.ctx.Done():
			continue
		default:
			return false
		}
	}
	return true
}

// Sends err to all receivers of endpoint and closes the channels.
func (c *client) sendErrToReceivers(err error, endpoint string) {
	result := requestResult{
		resp: nil,
		err:  err,
	}
	for _, recvr := range c.receivers[endpoint] {
		recvr.chn <- result
		close(recvr.chn)
	}
}

func closeResponseBody(b io.ReadCloser) {
	_ = b.Close()
}
