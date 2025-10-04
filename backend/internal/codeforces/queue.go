package codeforces

import "errors"

var ErrQueueEmpty = errors.New("queue is empty")

// Queue implementation using basic linked list.
type reqQueue struct {
	first *reqNode
	last  *reqNode
}

type reqNode struct {
	next *reqNode
	req  string
}

func (rq *reqQueue) push(s string) {
	n := &reqNode{
		next: nil,
		req:  s,
	}

	if rq.last == nil {
		rq.first = n
		rq.last = n
	} else {
		rq.last.next = n
		rq.last = n
	}
}

func (rq *reqQueue) front() (string, error) {
	if rq.first == nil {
		return "", ErrQueueEmpty
	}
	return rq.first.req, nil
}

func (rq *reqQueue) pop() error {
	if rq.first == nil {
		return ErrQueueEmpty
	}
	rq.first = rq.first.next
	return nil
}
