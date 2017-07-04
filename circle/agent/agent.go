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

// Command exports the agent command.
var Command = cli.Command{
	Name:   "agent",
	Usage:  "starts the circle agent",
	Action: loop,
}

func loop(c *cli.Context) error {
	endpoint := "ws://localhost:8000/ws/broker"
	filter := rpc.Filter{
		Labels: map[string]string{
			"platform": "linux",
		},
	}
	client, err := rpc.NewClient(
		endpoint,
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
	parallel := 1
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
	time.Sleep(time.Second * 1)

	path := "/Users/simon/Code/go/src/github.com/SimonXming/circle/test/pipeline-example.yaml"

	e := testRunPipeline(ctx, client, path)
	if e != nil {
		println(e)
	}

	// work, err := client.Next(ctx, filter)
	// if err != nil {
	// 	return err
	// }
	// println(work)
	os.Exit(1)
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
