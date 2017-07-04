package rpc2

import (
	"context"

	"github.com/SimonXming/pipeline/pipeline/backend"
)

type (
	Pipeline struct {
		ID      string          `json:"id"`
		Config  *backend.Config `json:"config"`
		Timeout int64           `json:"timeout"`
	}

	Filter struct {
		Labels map[string]string `json:"labels"`
		Expr   string            `json:"expr"`
	}
)

type Peer interface {
	// Next returns the next pipeline in the queue.
	Next(c context.Context, f Filter) (*Pipeline, error)
}
