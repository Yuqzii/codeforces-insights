package stats

import (
	"context"

	"github.com/yuqzii/cf-stats/internal/codeforces"
)

type Client interface {
	GetSubmissions(context.Context, string) ([]codeforces.Submission, error)
}

type service struct {
	client Client
}

func NewService(client Client) *service {
	return &service{
		client: client,
	}
}
