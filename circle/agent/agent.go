package agent

import (
	"context"
	"github.com/SimonXming/pipeline/pipeline/interrupt"
	// "github.com/SimonXming/pipeline/pipeline/rpc"
	rpc "github.com/SimonXming/pipeline/pipeline/rpc"
	"github.com/tevino/abool"
	"github.com/urfave/cli"
	"log"
	"math"
	"net/url"
	"sync"
	"time"
)

import (
	"github.com/SimonXming/pipeline/pipeline"
	"github.com/SimonXming/pipeline/pipeline/backend"
	"github.com/SimonXming/pipeline/pipeline/backend/docker"
	"github.com/SimonXming/pipeline/pipeline/multipart"
	"io"
	"os"
)

const (
	maxFileUpload = 5000000
	maxLogsUpload = 5000000
	maxProcs      = 1
	retryLimit    = math.MaxInt32
)

// Command exports the agent command.
var Command = cli.Command{
	Name:   "agent",
	Usage:  "starts the circle agent",
	Action: loop,
}

func loop(c *cli.Context) error {
	endpoint, err := url.Parse(
		"ws://localhost:8000/ws/broker",
	)
	if err != nil {
		return err
	}
	filter := rpc.Filter{
		Labels: map[string]string{
			"platform": "linux/amd64",
		},
	}
	client, err := rpc.NewClient(
		endpoint.String(),
		rpc.WithRetryLimit(
			retryLimit,
		),
		rpc.WithBackoff(
			time.Second*15,
		),
	)
	if err != nil {
		return err
	}
	defer client.Close()

	sigterm := abool.New()
	ctx := context.Background()
	ctx = interrupt.WithContextFunc(ctx, func() {
		println("ctrl+c received, terminating process")
		sigterm.Set()
	})

	var wg sync.WaitGroup
	parallel := maxProcs
	wg.Add(parallel)

	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			for {
				if sigterm.IsSet() {
					return
				}
				if err := run(ctx, client, filter); err != nil {
					log.Printf("build runner encountered error: exiting: %s", err)
					return
				}
			}
		}()
	}
	wg.Wait()
	return nil
}

/*
run 方法是 agent 的主要运行逻辑
1. 获取一个 job
2. 创建一个 docker engine
3. 处理等待 job 完成的逻辑(正确或错误)
4. 初始化 job
5. 给本次 job 设置 logger 和 tracer
6. 根据这次 job 的配置信息初始化 pipeline
7. 运行 pipeline 并实时更新 pipeline 状态
8. 完成 pipeline
9. 通过 connection 同步 pipelone 状态
*/

func run(ctx context.Context, client rpc.Peer, filter rpc.Filter) error {
	log.Println("pipeline: request next execution")

	// get the next job from the queue
	work, err := client.Next(ctx, filter)
	if err != nil {
		return err
	}
	if work == nil {
		return nil
	}
	log.Printf("pipeline: received next execution: %s", work.ID)

	// new docker engine
	engine, err := docker.NewEnv()
	if err != nil {
		return err
	}

	timeout := time.Hour
	if minutes := work.Timeout; minutes != 0 {
		timeout = time.Duration(minutes) * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cancelled := abool.New()
	go func() {
		if werr := client.Wait(ctx, work.ID); werr != nil {
			cancelled.SetTo(true)
			log.Printf("pipeline: cancel signal received: %s: %s", work.ID, werr)
			cancel()
		} else {
			log.Printf("pipeline: cancel channel closed: %s", work.ID)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("pipeline: cancel ping loop: %s", work.ID)
				return
			case <-time.After(time.Minute):
				log.Printf("pipeline: ping queue: %s", work.ID)
				client.Extend(ctx, work.ID)
			}
		}
	}()

	state := rpc.State{}
	state.Started = time.Now().Unix()
	err = client.Init(context.Background(), work.ID, state)
	if err != nil {
		log.Printf("pipeline: error signaling pipeline init: %s: %s", work.ID, err)
	}

	defaultLogger := pipeline.LogFunc(func(proc *backend.Step, rc multipart.Reader) error {
		part, err := rc.NextPart()
		if err != nil {
			return err
		}
		io.Copy(os.Stderr, part)
		return nil
	})

	err = pipeline.New(work.Config,
		pipeline.WithContext(ctx),
		pipeline.WithLogger(defaultLogger),
		// pipeline.WithTracer(defaultTracer),
		pipeline.WithEngine(engine),
	).Run()
	state.Finished = time.Now().Unix()
	state.Exited = true
	if err != nil {
		switch xerr := err.(type) {
		case *pipeline.ExitError:
			state.ExitCode = xerr.Code
		default:
			state.ExitCode = 1
			state.Error = err.Error()
		}
		if cancelled.IsSet() {
			state.ExitCode = 137
		}
	}

	log.Printf("pipeline: execution complete: %s", work.ID)

	err = client.Done(context.Background(), work.ID, state)
	if err != nil {
		log.Printf("Pipeine: error signaling pipeline done: %s: %s", work.ID, err)
	}

	return nil
}
