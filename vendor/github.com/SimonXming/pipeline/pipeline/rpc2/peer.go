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

	// State defines the pipeline state.
	State struct {
		Proc     string `json:"proc"`
		Exited   bool   `json:"exited"`
		ExitCode int    `json:"exit_code"`
		Started  int64  `json:"started"`
		Finished int64  `json:"finished"`
		Error    string `json:"error"`
	}

	Filter struct {
		Labels map[string]string `json:"labels"`
		Expr   string            `json:"expr"`
	}
)

type Peer interface {
	// Next returns the next pipeline in the queue.
	Next(c context.Context, f Filter) (*Pipeline, error)

	// Wait blocks until the pipeline is complete.
	Wait(c context.Context, id string) error

	// Init signals the pipeline is initialized.
	Init(c context.Context, id string, state State) error

	// Done signals the pipeline is complete.
	Done(c context.Context, id string, state State) error

	// Extend extends the pipeline deadline
	Extend(c context.Context, id string) error
}
