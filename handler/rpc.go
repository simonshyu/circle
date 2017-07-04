package handler

import (
	"context"
	"fmt"
	pipelineLib "github.com/SimonXming/pipeline/pipeline"
	"github.com/SimonXming/pipeline/pipeline/rpc2"
	"github.com/labstack/echo"
)

import (
	"io"
	"os"
)

type RPC struct {
	// queue queue.Queue
	queue []string
}

func RPCHandler(c echo.Context) error {
	temp_queue := make([]string, 0)
	temp_queue = append(temp_queue, "abc")
	peer := RPC{
		queue: temp_queue,
	}
	rpc2.NewServer(&peer).ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

// Next implements the rpc.Next function
func (s *RPC) Next(c context.Context, filter rpc2.Filter) (*rpc2.Pipeline, error) {
	fmt.Println(filter)

	pipeline := new(rpc2.Pipeline)

	path := "/Users/simon/Code/go/src/github.com/SimonXming/circle/test/task_data.json"
	var reader io.ReadCloser
	reader, err := os.Open(path)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}

	defer reader.Close()

	config, err := pipelineLib.Parse(reader)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	pipeline.ID = "1"
	pipeline.Config = config
	pipeline.Timeout = 10

	fmt.Printf("%v", pipeline)
	return pipeline, err
}
