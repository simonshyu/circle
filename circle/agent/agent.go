package agent

import (
	"context"
	"github.com/SimonXming/pipeline/pipeline/interrupt"
	"github.com/SimonXming/pipeline/pipeline/rpc"
	"github.com/SimonXming/pipeline/pipeline/rpc2"
	"github.com/tevino/abool"
	"github.com/urfave/cli"
	"log"
	"math"
	"net/url"
	"sync"
	"time"
)

import (
	"encoding/json"
	"fmt"

	"github.com/SimonXming/pipeline/pipeline"
	"github.com/SimonXming/pipeline/pipeline/backend"
	"github.com/SimonXming/pipeline/pipeline/backend/docker"
	"github.com/SimonXming/pipeline/pipeline/frontend"
	"github.com/SimonXming/pipeline/pipeline/frontend/yaml"
	"github.com/SimonXming/pipeline/pipeline/frontend/yaml/compiler"
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
	filter := rpc2.Filter{
		Labels: map[string]string{
			"platform": "linux/amd64",
		},
	}
	client, err := rpc2.NewClient(
		endpoint.String(),
		rpc2.WithRetryLimit(
			retryLimit,
		),
		rpc2.WithBackoff(
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

func run(ctx context.Context, client rpc2.Peer, filter rpc2.Filter) error {
	log.Println("pipeline: request next execution")

	// path := "/Users/simon/Code/go/src/github.com/SimonXming/circle/test/pipeline-example.yaml"

	// e := testRunPipeline(ctx, client, path)
	// if e != nil {
	// 	println(e)
	// }

	engine, err := docker.NewEnv()
	if err != nil {
		return err
	}

	defaultLogger := pipeline.LogFunc(func(proc *backend.Step, rc multipart.Reader) error {
		part, err := rc.NextPart()
		if err != nil {
			return err
		}
		io.Copy(os.Stderr, part)
		return nil
	})

	log.Printf("Trying get a work ...")
	work, err := client.Next(ctx, filter)
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

	state := rpc2.State{}
	state.Started = time.Now().Unix()
	err = client.Init(context.Background(), work.ID, state)
	if err != nil {
		log.Printf("pipeline: error signaling pipeline init: %s: %s", work.ID, err)
	}

	log.Printf("Success get a work.")
	err = pipeline.New(work.Config,
		pipeline.WithContext(ctx),
		pipeline.WithLogger(defaultLogger),
		// pipeline.WithTracer(defaultTracer),
		pipeline.WithEngine(engine),
	).Run()
	if err != nil {
		return err
	}

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

	log.Printf("pipeline: execution complete: %s", work.ID)
	err = client.Done(context.Background(), work.ID, state)
	if err != nil {
		log.Printf("Pipeine: error signaling pipeline done: %s: %s", work.ID, err)
	}

	return nil
}

func testRunPipeline(ctx context.Context, client rpc.Peer, filePath string) error {
	conf, err := yaml.ParseFile(filePath)
	if err != nil {
		return err
	}
	// fmt.Printf("%v", conf)
	compiled := compiler.New(
		compiler.WithPrefix(
			"test",
		),
		compiler.WithLocal(
			false,
		),
		compiler.WithMetadata(
			metadataFromContext(),
		),
		compiler.WithNetrc(
			"simon_xu@outlook.com",
			"@git5508177QaZ",
			"github.com",
		),
	).Compile(conf)

	err = outputJson(compiled)
	if err != nil {
		return err
	}

	engine, err := docker.NewEnv()
	if err != nil {
		return err
	}

	defaultLogger := pipeline.LogFunc(func(proc *backend.Step, rc multipart.Reader) error {
		part, err := rc.NextPart()
		if err != nil {
			return err
		}
		io.Copy(os.Stderr, part)
		return nil
	})

	err = pipeline.New(compiled,
		pipeline.WithContext(ctx),
		pipeline.WithLogger(defaultLogger),
		// pipeline.WithTracer(defaultTracer),
		pipeline.WithEngine(engine),
	).Run()

	if err != nil {
		fmt.Printf("%v", err)
	}
	return nil
}

func outputJson(compiled *backend.Config) error {
	if false {
		for _, stage := range compiled.Stages {
			for _, step := range stage.Steps {
				fmt.Printf("%v", step)
			}
		}
	}
	// marshal the compiled spec to formatted yaml
	out, err := json.MarshalIndent(compiled, "", "  ")
	if err != nil {
		return err
	}

	// create output file with option to dump to stdout
	var writer = os.Stdout
	output := "/Users/simon/Code/go/src/github.com/SimonXming/circle/test/pipeline-example.json"
	if output != "-" {
		writer, err = os.Create(output)
		if err != nil {
			return err
		}
	}
	defer writer.Close()

	_, err = writer.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func metadataFromContext() frontend.Metadata {
	return frontend.Metadata{
		Repo: frontend.Repo{
			Name:    "go-practice",
			Link:    "https://github.com/SimonXming/go-practice.git",
			Remote:  "https://github.com/SimonXming/go-practice.git",
			Private: false,
		},
		Curr: frontend.Build{
			Number:   1,
			Created:  0,
			Started:  0,
			Finished: 0,
			Status:   "start",
			Event:    "",
			Link:     "",
			Target:   "",
			Commit: frontend.Commit{
				Sha:     "50616752e10380848631c7c5bbabc87adb096d12",
				Ref:     "refs/heads/master",
				Refspec: "refs/heads/master",
				Branch:  "master",
				Message: "",
				Author: frontend.Author{
					Name:   "",
					Email:  "",
					Avatar: "",
				},
			},
		},
	}
}
