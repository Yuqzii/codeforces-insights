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
	ErrRateLimited            = errors.New("encountered Codeforces rate limit")
)

const requestBufferSize int = 1000

type client struct {
	client *http.Client
	url    string

	baseInterval time.Duration
	curInterval  time.Duration
	maxInterval  time.Duration
	muThrottling sync.RWMutex

	requests  chan string
	mu        sync.Mutex
	receivers map[string][]receiver
}

type clientOption func(*client)

func NewClient(httpClient *http.Client, url string, opts ...clientOption) *client {
	const (
		defaultBaseInterval = 2 * time.Second
		defaultMaxInterval  = 10 * time.Second
	)

	c := &client{
		client:       httpClient,
		url:          url,
		baseInterval: defaultBaseInterval,
		maxInterval:  defaultMaxInterval,
		requests:     make(chan string, requestBufferSize),
		receivers:    make(map[string][]receiver),
	}

	for _, opt := range opts {
		opt(c)
	}

	c.curInterval = c.baseInterval

	go c.listenForRequests()
	return c
}

func WithIntervals(baseInterval, maxInterval time.Duration) clientOption {
	return func(c *client) {
		c.baseInterval = baseInterval
		c.maxInterval = maxInterval
	}
}

type receiver struct {
	ctx context.Context
	chn chan<- requestResult
}

type requestResult struct {
	body []byte
	err  error
}

type apiResponse[T any] struct {
	Status  string `json:"status"`
	Result  []T    `json:"result"`
	Comment string `json:"comment,omitempty"`
}

// Adds the request to the queue. If the request already exists adds the caller as a receiver.
func (c *client) makeRequest(ctx context.Context, endpoint string) (<-chan requestResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

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
		return respChan, nil
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

	return respChan, nil
}

func (c *client) listenForRequests() {
	for {
		endpoint := <-c.requests

		if c.receiversCancelled(endpoint) {
			continue
		}

		err := c.sendRequest(endpoint)
		if err != nil {
			log.Printf("Error making request: %v\n", err)

			if errors.Is(err, ErrRateLimited) {
				c.adjustThrottle(2) // Double interval.
			}
		} else {
			c.adjustThrottle(0.8)
		}

		c.muThrottling.RLock()
		interval := c.curInterval
		c.muThrottling.RUnlock()

		time.Sleep(interval)
	}
}

// Sends the request for the specified endpoint and broadcasts the result to all receivers.
func (c *client) sendRequest(endpoint string) error {
	resp, err := c.client.Get(c.url + endpoint)
	if err != nil {
		c.sendErrToReceivers(err, endpoint)
		return fmt.Errorf("requesting '%s' from Codeforces: %w", endpoint, err)
	}

	if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusTooManyRequests {
		// Likely being rate limited by Codeforces.
		// They usually return a 405 instead of the arguable more correct 429 for rate limiting.
		err = fmt.Errorf("%w: %s", ErrRateLimited, resp.Status)
		c.sendErrToReceivers(err, endpoint)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close() // nolint:errcheck
	if err != nil {
		c.sendErrToReceivers(err, endpoint)
		return fmt.Errorf("reading '%s' response body: %w", endpoint, err)
	}

	result := requestResult{
		body: body,
		err:  nil,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, recvr := range c.receivers[endpoint] {
		select {
		case <-recvr.ctx.Done(): // Don't send to cancelled receiver
		default:
			recvr.chn <- result
		}
		close(recvr.chn)
	}

	delete(c.receivers, endpoint)

	return nil
}

func (c *client) adjustThrottle(factor float64) {
	c.muThrottling.Lock()

	newInterval := time.Duration(float64(c.curInterval) * factor)

	// Clamp interval.
	c.curInterval = min(newInterval, c.maxInterval)
	c.curInterval = max(c.curInterval, c.baseInterval)

	c.muThrottling.Unlock()
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
		body: nil,
		err:  err,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, recvr := range c.receivers[endpoint] {
		select {
		case <-recvr.ctx.Done(): // Don't send to cancelled receiver
		default:
			recvr.chn <- result
		}
		close(recvr.chn)
	}
}
