package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SimonXming/circle/model"
	pipelineLib "github.com/SimonXming/pipeline/pipeline"
	"github.com/SimonXming/pipeline/pipeline/rpc2"
	"github.com/SimonXming/queue"
	"github.com/labstack/echo"
)

import (
	"io"
	"io/ioutil"
	"os"
)

var Config = struct {
	Services struct {
		// Pubsub     pubsub.Publisher
		Queue queue.Queue
		// Logs       logging.Log
		// Senders    model.SenderService
		// Secrets    model.SecretService
		// Registries model.RegistryService
		// Environ    model.EnvironService
	}
	Storage struct {
		// Users  model.UserStore
		// Repos  model.RepoStore
		// Builds model.BuildStore
		// Logs   model.LogStore
		Config model.ConfigStore
		// Registries model.RegistryStore
		// Secrets model.SecretStore
	}
	Pipeline struct {
		Limits     model.ResourceLimit
		Volumes    []string
		Networks   []string
		Privileged []string
	}
}{}

type RPC struct {
	queue queue.Queue
	// queue []string
}

func RPCHandler(c echo.Context) error {
	// temp_queue := make([]string, 0)
	// temp_queue = append(temp_queue, "abc")
	peer := RPC{
		queue: Config.Services.Queue,
	}
	server := rpc2.NewServer(&peer)
	server.ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

// Next implements the rpc.Next function
func (s *RPC) NextBak(c context.Context, filter rpc2.Filter) (*rpc2.Pipeline, error) {
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

// Next implements the rpc.Next function
func (s *RPC) Next(c context.Context, filter rpc2.Filter) (*rpc2.Pipeline, error) {
	fn := func(task *queue.Task) bool {
		for k, v := range filter.Labels {
			if task.Labels[k] != v {
				return false
			}
		}
		return true
	}

	task, err := s.queue.Poll(c, fn)
	if err != nil {
		return nil, err
	} else if task == nil {
		return nil, nil
	}
	pipeline := new(rpc2.Pipeline)

	// check if the process was previously cancelled
	// cancelled, _ := s.checkCancelled(pipeline)
	// if cancelled {
	// 	logrus.Debugf("ignore pid %v: cancelled by user", pipeline.ID)
	// 	if derr := s.queue.Done(c, pipeline.ID); derr != nil {
	// 		logrus.Errorf("error: done: cannot ack proc_id %v: %s", pipeline.ID, err)
	// 	}
	// 	return nil, nil
	// }

	err = json.Unmarshal(task.Data, pipeline)

	path := "/Users/simon/Code/go/src/github.com/SimonXming/circle/test/tmp_task_data.json"

	pipelineJson, _ := json.Marshal(pipeline)
	err = ioutil.WriteFile(path, pipelineJson, 0644)

	return pipeline, err
}
