package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"math/rand"
	"time"

	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/remote"
	"github.com/SimonXming/circle/store"
	"github.com/SimonXming/circle/utils/httputil"
	"github.com/SimonXming/pipeline/pipeline/backend"
	"github.com/SimonXming/pipeline/pipeline/frontend"
	"github.com/SimonXming/pipeline/pipeline/frontend/yaml"
	"github.com/SimonXming/pipeline/pipeline/frontend/yaml/compiler"
	"github.com/SimonXming/pipeline/pipeline/frontend/yaml/matrix"
	"github.com/SimonXming/pipeline/pipeline/rpc2"
	"github.com/SimonXming/queue"

	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func PostBuild(c echo.Context) error {
	cc, ok := c.(*CircleContext)
	if !ok {
		c.String(http.StatusBadRequest, "Context 转换失败.")
		return errors.New("Context 转换失败.")
	}
	_type := cc.DefaultQueryParam("type", model.BuildManual)
	repoID, err := strconv.ParseInt(c.Param("repoID"), 10, 64)
	if err != nil {
		c.Error(err)
		return err
	}

	repo, err := store.RepoLoad(c, repoID)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	conf, err := store.ConfigFind(c, repo)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return err
	}

	account, err := store.SetupRemoteWithScmID(c, repo.ScmId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	netrc, err := remote.Netrc(c, account)
	if err != nil {
		logrus.Errorf("failure to generate netrc for %s. %s", repo.FullName, err)
		c.String(http.StatusNotFound, err.Error())
		return err
	}

	var build *model.Build

	if _type == model.BuildManual {
		build = new(model.Build)
		build.RepoID = repoID
		build.Number = 0
		build.Event = model.EventManual
		build.Ref = fmt.Sprintf("refs/heads/%s", repo.Branch)
		build.Branch = repo.Branch
		build.Refspec = fmt.Sprintf("refs/heads/%s", repo.Branch)
		build.Status = model.StatusPending
		build.Started = 0
		build.Finished = 0
		build.Enqueued = time.Now().UTC().Unix()
		build.Error = ""
		err = store.BuildCreate(c, build)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return err
		}
	} else if _type == model.BuildFork {
		num, err := strconv.Atoi(c.QueryParam("number"))
		if err != nil {
			c.Error(err)
			return err
		}
		build, err = store.GetBuildNumber(c, repo, num)
		if err != nil {
			logrus.Errorf("failure to get build %d. %s", num, err)
			c.String(http.StatusNotFound, err.Error())
			return err
		}
		build.ID = 0
		build.Number = 0
		// build.Event: using the same event
		build.Status = model.StatusPending
		build.Started = 0
		build.Finished = 0
		build.Enqueued = time.Now().UTC().Unix()
		build.Error = ""
		err = store.BuildCreate(c, build)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return err
		}
	} else if _type == model.BuildRerun {
		num, err := strconv.Atoi(c.QueryParam("number"))
		if err != nil {
			c.Error(err)
			return err
		}
		build, err = store.GetBuildNumber(c, repo, num)
		if err != nil {
			logrus.Errorf("failure to get build %d. %s", num, err)
			c.String(http.StatusNotFound, err.Error())
			return err
		}
		build.Status = model.StatusPending
		build.Started = 0
		build.Finished = 0
		build.Enqueued = time.Now().UTC().Unix()
		build.Error = ""
		err = store.ProcClear(c, build)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return err
		}

		err = store.BuildUpdate(c, build)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return err
		}
	}

	envs := map[string]string{}
	envs["someEnv"] = "someValue"

	b := builder{
		Repo:  repo,
		Curr:  build,
		Netrc: netrc,
		Envs:  envs,
		Link:  httputil.GetURL(c.Request()),
		Yaml:  conf.Data,
	}
	items, err := b.Build()
	if err != nil {
		build.Status = model.StatusError
		build.Started = time.Now().Unix()
		build.Finished = build.Started
		build.Error = err.Error()
		store.BuildUpdate(c, build)
		return err
	}

	// 如果没有用到 matrix 语法，
	// len(items) == 1

	// pcounter 代表构建流程的个数
	var pcounter = len(items)
	for _, item := range items {
		// initProc := *(item.Proc)
		build.Procs = append(build.Procs, item.Proc)
		item.Proc.BuildID = build.ID

		for _, stage := range item.Config.Stages {
			// 初始化 gid 的值
			var gid int
			for _, step := range stage.Steps {
				// 每多一个 step，pcounter 就会自增 +1
				// ** 因此 PID 用于标示 step
				pcounter++
				if gid == 0 {
					// 在每个 stage 的第一个 step 时令 gid = pcounter
					// ** PGID 用于标示 stage 的值
					// 第一个 stage 的gid值是 （matrix 个数 +1）
					// 第二个 stage 的gid值是 （matrix 个数 + 上一个stage的step数 + 1）
					// 第 n 个 stage 的gid值是 （matrix 个数 + 第 n-1 个stage的step数 + 1）
					gid = pcounter
				}
				// ** PPID 用于标示 matrix item
				proc := &model.Proc{
					BuildID: build.ID,
					Name:    step.Alias,
					PID:     pcounter,
					PPID:    item.Proc.PID,
					PGID:    gid,
					State:   model.StatusPending,
				}
				build.Procs = append(build.Procs, proc)
			}
		}
	}

	err = store.ProcCreate(c, build.Procs)
	if err != nil {
		logrus.Errorf("error persisting procs %s/%d: %s", repo.FullName, build.Number, err)
	}

	c.JSON(http.StatusAccepted, build)

	// 如果没有用到 matrix 语法，
	// len(items) == 1
	for _, item := range items {
		task := new(queue.Task)
		task.ID = fmt.Sprint(item.Proc.ID)
		task.Labels = map[string]string{}
		task.Labels["platform"] = item.Platform
		for k, v := range item.Labels {
			task.Labels[k] = v
		}

		task.Data, _ = json.Marshal(rpc2.Pipeline{
			ID:      fmt.Sprint(item.Proc.ID),
			Config:  item.Config,
			Timeout: b.Repo.Timeout,
		})

		// Config.Services.Logs.Open(context.Background(), task.ID)
		Config.Services.Queue.Push(context.Background(), task)
	}

	return nil
}

// return the metadata from the cli context.
func metadataFromStruct(repo *model.Repo, build *model.Build, proc *model.Proc, link string) frontend.Metadata {
	return frontend.Metadata{
		Repo: frontend.Repo{
			Name:   repo.FullName,
			Link:   repo.Link,
			Remote: repo.Clone,
			// Private: repo.IsPrivate,
		},
		Curr: frontend.Build{
			Number:   build.Number,
			Created:  build.Created,
			Started:  build.Started,
			Finished: build.Finished,
			Status:   build.Status,
			Event:    build.Event,
			Link:     build.Link,
			Commit: frontend.Commit{
				Sha:     build.Commit,
				Ref:     build.Ref,
				Refspec: build.Refspec,
				Branch:  build.Branch,
			},
		},
		Job: frontend.Job{
			Number: proc.PID,
			Matrix: proc.Environ,
		},
		Sys: frontend.System{
			Name: "circle",
			Link: link,
			Arch: "linux/amd64",
		},
	}
}

type builder struct {
	Repo *model.Repo
	Curr *model.Build
	// Last  *model.Build
	Netrc *model.Netrc
	// Secs []*model.Secret
	// Regs []*model.Registry
	Link string
	Yaml string
	Envs map[string]string
}

type buildItem struct {
	Proc     *model.Proc
	Platform string
	Labels   map[string]string
	Config   *backend.Config
}

func (b *builder) Build() ([]*buildItem, error) {
	axes, err := matrix.ParseString(b.Yaml)
	if err != nil {
		return nil, err
	}
	if len(axes) == 0 {
		axes = append(axes, matrix.Axis{})
	}
	var items []*buildItem
	for i, axis := range axes {
		proc := &model.Proc{
			BuildID: b.Curr.ID,
			PID:     i + 1,
			PGID:    i + 1,
			State:   model.StatusPending,
			Environ: axis,
		}

		// 定义每一步所需要的环境变量
		metadata := metadataFromStruct(b.Repo, b.Curr, proc, b.Link)
		environ := metadata.Environ()

		for k, v := range axis {
			environ[k] = v
		}

		parsed, err := yaml.ParseString(b.Yaml)
		if err != nil {
			return nil, err
		}
		metadata.Sys.Arch = parsed.Platform
		if metadata.Sys.Arch == "" {
			metadata.Sys.Arch = "linux/amd64"
		}

		ir := compiler.New(
			compiler.WithEnviron(environ),
			compiler.WithEnviron(b.Envs),
			compiler.WithEscalated(Config.Pipeline.Privileged...),
			compiler.WithResourceLimit(Config.Pipeline.Limits.MemSwapLimit, Config.Pipeline.Limits.MemLimit, Config.Pipeline.Limits.ShmSize, Config.Pipeline.Limits.CPUQuota, Config.Pipeline.Limits.CPUShares, Config.Pipeline.Limits.CPUSet),
			compiler.WithVolumes(Config.Pipeline.Volumes...),
			compiler.WithNetworks(Config.Pipeline.Networks...),
			compiler.WithLocal(false),
			compiler.WithOption(
				compiler.WithNetrc(
					b.Netrc.Login,
					b.Netrc.Password,
					b.Netrc.Machine,
				),
				b.Repo.IsPrivate,
			),
			// compiler.WithRegistry(registries...),
			// compiler.WithSecret(secrets...),
			compiler.WithPrefix(
				fmt.Sprintf(
					"%d_%d",
					proc.ID,
					rand.Int(),
				),
			),
			compiler.WithEnviron(proc.Environ),
			compiler.WithProxy(),
			compiler.WithWorkspaceFromURL("/circle", b.Curr.Link),
			compiler.WithMetadata(metadata),
		).Compile(parsed)

		item := &buildItem{
			Proc:     proc,
			Config:   ir,
			Labels:   parsed.Labels,
			Platform: metadata.Sys.Arch,
		}
		if item.Labels == nil {
			item.Labels = map[string]string{}
		}
		items = append(items, item)
	}
	return items, nil

}
