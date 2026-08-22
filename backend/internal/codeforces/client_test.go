package codeforces

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCancelledRequestCanBeQueuedAgain(t *testing.T) {
	c := &client{
		requests:  make(chan string, requestBufferSize),
		receivers: make(map[string][]receiver),
	}
	const endpoint = "problemset.problems"

	ctx, cancel := context.WithCancel(context.Background())
	_, err := c.makeRequest(ctx, endpoint)
	require.NoError(t, err)
	require.Equal(t, endpoint, <-c.requests)

	cancel()
	require.True(t, c.receiversCancelled(endpoint))

	_, err = c.makeRequest(context.Background(), endpoint)
	require.NoError(t, err)

	select {
	case got := <-c.requests:
		require.Equal(t, endpoint, got)
	default:
		t.Fatal("request was not queued again after its previous receiver cancelled")
	}
}
