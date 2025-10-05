package codeforces

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	ErrCodeforcesReturnedFail = errors.New("Codeforces returned status FAILED") // nolint:staticcheck
	ErrAllReceiversCancelled  = errors.New("all receivers to request cancelled")
)

const requestBufferSize int = 1000

type client struct {
	client          *http.Client
	url             string
	timeBetweenReqs time.Duration

	requests  chan string
	mu        sync.Mutex
	receivers map[string][]receiver
}

func NewClient(httpClient *http.Client, url string, timeBetweenReqs time.Duration) *client {
	c := &client{
		client:          httpClient,
		url:             url,
		timeBetweenReqs: timeBetweenReqs,
		requests:        make(chan string, requestBufferSize),
		receivers:       make(map[string][]receiver),
	}
	c.listenForRequests()
	return c
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

func (c *client) makeRequest(ctx context.Context, endpoint string) <-chan requestResult {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, queued := c.receivers[endpoint]
	if queued {
		// Request is already queued, just add to receivers
		respChan := make(chan requestResult)
		c.receivers[endpoint] = append(c.receivers[endpoint], receiver{
			ctx: ctx,
			chn: respChan,
		})
		return respChan
	}

	// Create receiver list
	c.receivers[endpoint] = make([]receiver, 0, 1)

	respChan := make(chan requestResult)
	c.receivers[endpoint] = append(c.receivers[endpoint], receiver{
		ctx: ctx,
		chn: respChan,
	})

	// Push request to queue
	c.requests <- endpoint

	return respChan
}

func (c *client) listenForRequests() {
	go func() {
		for {
			endpoint := <-c.requests

			if c.receiversCancelled(endpoint) {
				continue
			}

			t := time.Now()
			err := c.sendRequest(endpoint)
			if err != nil {
				log.Printf("Error sending request: %v\n", err)
			}
			time.Sleep(c.timeBetweenReqs - time.Since(t))
		}
	}()
}

// Sends the request for the specified endpoint and broadcasts the result to all receivers.
func (c *client) sendRequest(endpoint string) error {
	resp, err := c.client.Get(c.url + endpoint)
	if err != nil {
		c.sendErrToReceivers(err, endpoint)
		return fmt.Errorf("requesting '%s' from Codeforces: %w", endpoint, err)
	}

	result := requestResult{
		resp: resp,
		err:  nil,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, recvr := range c.receivers[endpoint] {
		recvr.chn <- result
		close(recvr.chn)
	}

	delete(c.receivers, endpoint)

	return nil
}

// Returns true if all receivers to endpoint has cancelled their context.
func (c *client) receiversCancelled(endpoint string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, recvr := range c.receivers[endpoint] {
		recvr.chn <- result
		close(recvr.chn)
	}
}

func closeResponseBody(b io.ReadCloser) {
	_ = b.Close()
}
