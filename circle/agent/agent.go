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
	parallel := 5
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
	time.Sleep(time.Second * 5)
	work, err := client.Next(ctx, filter)
	if err != nil {
		return err
	}
	println(work)
	return nil
}
