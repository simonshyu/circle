package handler

// TODO: convert rpc2 to rpc

import (
	"context"
	"encoding/json"
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/store"
	"github.com/SimonXming/pipeline/pipeline/rpc2"
	"github.com/SimonXming/queue"
	"github.com/labstack/echo"
	"log"
	"strconv"
)

import (
	"io/ioutil"
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
	store store.Store
}

func RPCHandler(c echo.Context) error {
	peer := RPC{
		store: store.FromContext(c),
		queue: Config.Services.Queue,
	}
	server := rpc2.NewServer(&peer)
	server.ServeHTTP(c.Response().Writer, c.Request())
	return nil
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

	// Output next process workon in json
	path := "/Users/simon/Code/go/src/github.com/SimonXming/circle/test/next_task_data.json"
	pipelineJson, _ := json.Marshal(pipeline)
	err = ioutil.WriteFile(path, pipelineJson, 0644)

	return pipeline, err
}

// Init implements the rpc.Init function
func (s *RPC) Init(c context.Context, id string, state rpc2.State) error {
	procID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	proc, err := s.store.ProcLoad(procID)
	if err != nil {
		log.Printf("error: cannot find proc with id %d: %s", procID, err)
		return err
	}

	build, err := s.store.BuildLoad(proc.BuildID)
	if err != nil {
		log.Printf("error: cannot find build with id %d: %s", proc.BuildID, err)
		return err
	}

	// repo, err := s.store.RepoLoad(build.RepoID)
	// if err != nil {
	// 	log.Printf("error: cannot find repo with id %d: %s", build.RepoID, err)
	// 	return err
	// }

	if build.Status == model.StatusPending {
		build.Status = model.StatusRunning
		build.Started = state.Started
		if err := s.store.BuildUpdate(build); err != nil {
			log.Printf("error: init: cannot update build_id %d state: %s", build.ID, err)
		}
	}

	// defer func() {
	// 	build.Procs, _ = s.store.ProcList(build)
	// 	message := pubsub.Message{
	// 		Labels: map[string]string{
	// 			"repo":    repo.FullName,
	// 			"private": strconv.FormatBool(repo.IsPrivate),
	// 		},
	// 	}
	// 	message.Data, _ = json.Marshal(model.Event{
	// 		Repo:  *repo,
	// 		Build: *build,
	// 	})
	// 	s.pubsub.Publish(c, "topic/events", message)
	// }()

	proc.Started = state.Started
	proc.State = model.StatusRunning
	return s.store.ProcUpdate(proc)
}

// Done implements the rpc.Done function
func (s *RPC) Done(c context.Context, id string, state rpc2.State) error {
	procID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}

	proc, err := s.store.ProcLoad(procID)
	if err != nil {
		log.Printf("error: cannot find proc with id %d: %s", procID, err)
		return err
	}

	build, err := s.store.BuildLoad(proc.BuildID)
	if err != nil {
		log.Printf("error: cannot find build with id %d: %s", proc.BuildID, err)
		return err
	}

	// repo, err := s.store.RepoLoad(build.RepoID)
	// if err != nil {
	// 	log.Printf("error: cannot find repo with id %d: %s", build.RepoID, err)
	// 	return err
	// }

	proc.Stopped = state.Finished
	proc.Error = state.Error
	proc.ExitCode = state.ExitCode
	proc.State = model.StatusSuccess
	if proc.ExitCode != 0 || proc.Error != "" {
		proc.State = model.StatusFailure
	}
	if err := s.store.ProcUpdate(proc); err != nil {
		log.Printf("error: done: cannot update proc_id %d state: %s", procID, err)
	}

	if err := s.queue.Done(c, id); err != nil {
		log.Printf("error: done: cannot ack proc_id %d: %s", procID, err)
	}

	// TODO handle this error
	procs, _ := s.store.ProcList(build)
	for _, p := range procs {
		if p.Running() && p.PPID == proc.PID {
			p.State = model.StatusSkipped
			if p.Started != 0 {
				p.State = model.StatusSuccess // for deamons that are killed
				p.Stopped = proc.Stopped
			}
			if err := s.store.ProcUpdate(p); err != nil {
				log.Printf("error: done: cannot update proc_id %d child state: %s", p.ID, err)
			}
		}
	}

	running := false
	status := model.StatusSuccess
	for _, p := range procs {
		if p.PPID == 0 {
			if p.Running() {
				running = true
			}
			if p.Failing() {
				status = p.State
			}
		}
	}
	if !running {
		build.Status = status
		build.Finished = proc.Stopped
		if err := s.store.BuildUpdate(build); err != nil {
			log.Printf("error: done: cannot update build_id %d final state: %s", build.ID, err)
		}
	}

	// if err := s.logger.Close(c, id); err != nil {
	// 	log.Printf("error: done: cannot close build_id %d logger: %s", proc.ID, err)
	// }

	// build.Procs = model.Tree(procs)
	// message := pubsub.Message{
	// 	Labels: map[string]string{
	// 		"repo":    repo.FullName,
	// 		"private": strconv.FormatBool(repo.IsPrivate),
	// 	},
	// }
	// message.Data, _ = json.Marshal(model.Event{
	// 	Repo:  *repo,
	// 	Build: *build,
	// })
	// s.pubsub.Publish(c, "topic/events", message)

	return nil
}
