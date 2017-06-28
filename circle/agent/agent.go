package agent

import (
	"context"
	"github.com/SimonXming/pipeline/pipeline/interrupt"
	"github.com/SimonXming/pipeline/pipeline/rpc"
	"github.com/tevino/abool"
	"github.com/urfave/cli"
	"log"
	"sync"
	"time"
)

// Command exports the agent command.
var Command = cli.Command{
	Name:   "agent",
	Usage:  "starts the circle agent",
	Action: loop,
}

func loop(c *cli.Context) error {
	endpoint := "ws://localhost:8000/ws/broker"
	client, err := rpc.NewClient(
		endpoint,
	)

	sigterm := abool.New()
	ctx := context.Background()
	ctx = interrupt.WithContextFunc(ctx, func() {
		println("ctrl+c received, terminating process")
		sigterm.Set()
	})

	var wg sync.WaitGroup
	parallel := 5
	wg.Add(parallel)

	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			for {
				if sigterm.IsSet() {
					return
				}
				if err := run(ctx); err != nil {
					log.Printf("build runner encountered error: exiting: %s", err)
					return
				}
			}
		}()
	}
	wg.Wait()
	return nil
}

func run(ctx context.Context) error {
	log.Println("pipeline: request next execution")
	time.Sleep(time.Second * 5)
	return nil
}
