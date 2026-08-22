package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yuqzii/codeforces-insights/internal/codeforces"
)

const (
	errorContestID    = 1
	lateContestID     = 2
	sentinelContestID = 3
)

type channelCloseContestResultsProvider struct {
	handleContext chan context.Context
	lateStarted   chan struct{}
	releaseLate   chan struct{}
}

func (p *channelCloseContestResultsProvider) GetContestResults(
	_ context.Context, id int,
) ([]codeforces.Contestant, *codeforces.Contest, error) {
	switch id {
	case errorContestID:
		<-p.lateStarted
		return nil, nil, errors.New("contest results error")
	case lateContestID:
		close(p.lateStarted)
		<-p.releaseLate
	}

	return []codeforces.Contestant{{OldRating: 1500}}, &codeforces.Contest{}, nil
}

func (p *channelCloseContestResultsProvider) GetContestResultsFromHandle(
	ctx context.Context, _ string,
) ([]codeforces.Contestant, error) {
	p.handleContext <- ctx
	return nil, nil
}

func TestHandlePerformanceLeavesResultChannelOpenForInFlightJobs(t *testing.T) {
	provider := &channelCloseContestResultsProvider{
		handleContext: make(chan context.Context, 1),
		lateStarted:   make(chan struct{}),
		releaseLate:   make(chan struct{}),
	}
	h := New(nil, provider, nil, nil, 3, 2)

	body := `{
		"handle":"tourist",
		"ratingHistory":[
			{"contestId":1,"rank":1,"oldRating":1500},
			{"contestId":2,"rank":1,"oldRating":1500}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/performance", strings.NewReader(body))
	res := httptest.NewRecorder()
	handlerDone := make(chan struct{})
	go func() {
		h.HandlePerformance(res, req)
		close(handlerDone)
	}()

	var handlerCtx context.Context
	select {
	case handlerCtx = <-provider.handleContext:
	case <-time.After(5 * time.Second):
		t.Fatal("handler did not request contest results")
	}

	select {
	case <-handlerDone:
	case <-time.After(5 * time.Second):
		t.Fatal("handler did not return after a performance job failed")
	}

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, res.Code)
	}
	select {
	case <-handlerCtx.Done():
	default:
		t.Fatal("handler context was not canceled after a performance job failed")
	}

	// Queue a sentinel behind the blocked job. Receiving its result proves that
	// the in-flight job completed and sent its result after handler cancellation.
	sentinelResult := make(chan perfResult, 1)
	h.perf.addJob(context.Background(), &codeforces.RatingChange{
		ContestID: sentinelContestID,
		Rank:      1,
		OldRating: 1500,
	}, sentinelResult)
	close(provider.releaseLate)

	select {
	case result := <-sentinelResult:
		if result.err != nil {
			t.Fatalf("sentinel job failed: %v", result.err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("in-flight worker did not continue after sending its late result")
	}
}
