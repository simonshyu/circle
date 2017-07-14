package handler

import (
	"encoding/base32"
	"fmt"
	"github.com/gorilla/securecookie"

	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/remote"
	"github.com/SimonXming/circle/store"
	"github.com/SimonXming/circle/utils/httputil"
	"github.com/SimonXming/circle/utils/token"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func GetAllRepo(c echo.Context) error {
	repos, err := store.RepoList(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	return c.JSON(http.StatusOK, repos)
}

func GetRepos(c echo.Context) error {
	scmId, err := strconv.ParseInt(c.Param("scmID"), 10, 64)
	if err != nil {
		c.Error(err)
		return err
	}
	account, err := store.ScmAccountLoad(c, scmId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	repos, err := store.RepoFind(c, account)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	return c.JSON(http.StatusOK, repos)
}

func GetRemoteRepos(c echo.Context) error {
	scmId, err := strconv.ParseInt(c.Param("scmID"), 10, 64)
	if err != nil {
		c.Error(err)
		return err
	}
	_, err = store.SetupRemoteWithScmID(c, scmId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	repos, err := remote.Repos(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, repos)
}

func PostRepo(c echo.Context) error {
	scmId, err := strconv.ParseInt(c.Param("scmID"), 10, 64)
	if err != nil {
		c.Error(err)
		return err
	}

	in := new(model.Repo)
	if err := c.Bind(in); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	_, err = store.SetupRemoteWithScmID(c, scmId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	r, err := remote.Repo(c, in.Owner, in.Name)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	// error if the repository already exists
	_, err = store.GetRepoScmIDOwnerName(c, scmId, r.Owner, r.Name)
	if err == nil {
		c.String(http.StatusConflict, "Repository already exists.")
		return err
	}
	r.ScmId = scmId
	r.AllowPush = true
	r.AllowPull = true
	r.Hash = base32.StdEncoding.EncodeToString(
		securecookie.GenerateRandomKey(32),
	)

	// crates the jwt token used to verify the repository
	t := token.New(token.HookToken, r.FullName)
	sig, err := t.Sign(r.Hash)
	if err != nil {
		c.String(http.StatusBadRequest, "Generate webhook token failed.")
		return err
	}

	link := fmt.Sprintf(
		"%s/hook?access_token=%s",
		httputil.GetURL(c.Request()),
		sig,
	)

	// activate the repository before we make any
	// local changes to the database.
	err = remote.Activate(c, r, link)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	err = store.RepoCreate(c, r)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	return nil
}
