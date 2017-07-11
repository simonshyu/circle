package handler

import (
	// "crypto/sha256"
	// "encoding/base32"
	"fmt"
	"time"
	// "github.com/gorilla/securecookie"

	"github.com/SimonXming/circle/model"
	// "github.com/SimonXming/circle/remote"
	"github.com/SimonXming/circle/store"
	// "github.com/SimonXming/circle/utils"
	"github.com/SimonXming/circle/utils/httputil"
	// "github.com/SimonXming/circle/utils/token"
	"github.com/SimonXming/pipeline/pipeline/frontend/yaml/matrix"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

type builder struct {
	Repo *model.Repo
	Curr *model.Build
	// Last  *model.Build
	// Netrc *model.Netrc
	// Secs []*model.Secret
	// Regs []*model.Registry
	Link string
	Yaml string
	Envs map[string]string
}

func PostBuild(c echo.Context) error {
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
	fmt.Printf("%v", conf.Data)

	build := new(model.Build)
	build.RepoID = repoID
	build.Number = 0
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

	envs := map[string]string{}
	envs["someEnv"] = "someValue"

	b := builder{
		Repo: repo,
		Curr: build,
		Envs: envs,
		Link: httputil.GetURL(c.Request()),
		Yaml: conf.Data,
	}
	b.Build()

	return nil
}

func (b *builder) Build() error {
	axes, err := matrix.ParseString(b.Yaml)
	if err != nil {
		return err
	}
	if len(axes) == 0 {
		axes = append(axes, matrix.Axis{})
	}

	return nil
	// for i, axis := range axes {
	// 	metadata := metadataFromStruct(b.Repo, b.Curr, b.Last, proc, b.Link)
	// }
}
