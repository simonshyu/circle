package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"time"

	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/remote"
	"github.com/SimonXming/circle/store"
	"github.com/SimonXming/circle/utils/httputil"
	"github.com/SimonXming/circle/utils/token"
	rpc "github.com/SimonXming/pipeline/pipeline/rpc2"
	"github.com/SimonXming/queue"

	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func PostHook(c echo.Context) error {
	scmId, err := strconv.ParseInt(c.QueryParam("scmID"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	account, err := store.SetupRemoteWithScmID(c, scmId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	tmprepo, build, err := remote.Hook(c, c.Request())
	if err != nil {
		logrus.Errorf("failure to parse hook. %s", err)
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	if build == nil {
		logrus.Errorf("Failed generate build from hook request.")
		c.String(http.StatusBadRequest, "Failed generate build from hook request.")
		return err
	}
	if tmprepo == nil {
		logrus.Errorf("failure to ascertain repo from hook.")
		c.String(http.StatusBadRequest, "failure to ascertain repo from hook.")
		return err
	}
	repo, err := store.GetRepoOwnerName(c, tmprepo.Owner, tmprepo.Name)
	if err != nil {
		logrus.Errorf("failure to find repo %s/%s from hook. %s", tmprepo.Owner, tmprepo.Name, err)
		c.String(http.StatusNotFound, "failure to ascertain repo from hook.")
		return err
	}

	// get the token and verify the hook is authorized
	parsed, err := token.ParseRequest(c.Request(), func(t *token.Token) (string, error) {
		return repo.Hash, nil
	})
	if err != nil {
		logrus.Errorf("failure to parse token from hook for %s. %s", repo.FullName, err)
		c.String(http.StatusBadRequest, "failure to ascertain repo from hook.")
		return err
	}
	if parsed.Text != repo.FullName {
		logrus.Errorf("failure to verify token from hook. Expected %s, got %s", repo.FullName, parsed.Text)
		c.String(http.StatusForbidden, "failure to ascertain repo from hook.")
		return err
	}
	var skipped = true
	if (build.Event == model.EventPush && repo.AllowPush) ||
		(build.Event == model.EventPull && repo.AllowPull) ||
		(build.Event == model.EventTag && repo.AllowTag) {
		skipped = false
	}

	if skipped {
		logrus.Infof("ignoring hook. repo %s is disabled for %s events.", repo.FullName, build.Event)
		c.String(http.StatusNoContent, "failure to ascertain repo from hook.")
		return nil
	}
	conf, err := store.ConfigFind(c, repo)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return err
	}
	build.ConfigID = conf.ID

	netrc, err := remote.Netrc(c, account)
	if err != nil {
		logrus.Errorf("failure to generate netrc for %s. %s", repo.FullName, err)
		c.String(http.StatusNotFound, err.Error())
		return err
	}

	// verify the branches can be built vs skipped
	// branches, err := yaml.ParseString(conf.Data)
	// if err == nil {
	// 	if !branches.Branches.Match(build.Branch) && build.Event != model.EventTag && build.Event != model.EventDeploy {
	// 		c.String(200, "Branch does not match restrictions defined in yaml")
	// 		return
	// 	}
	// }
	// update some build fields
	build.RepoID = repo.ID
	build.Status = model.StatusPending
	err = store.BuildCreate(c, build, build.Procs...)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
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

	c.JSON(http.StatusCreated, build)

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

		task.Data, _ = json.Marshal(rpc.Pipeline{
			ID:      fmt.Sprint(item.Proc.ID),
			Config:  item.Config,
			Timeout: b.Repo.Timeout,
		})

		// Config.Services.Logs.Open(context.Background(), task.ID)
		Config.Services.Queue.Push(context.Background(), task)
	}

	return nil
}
