package handler

import (
	// "fmt"
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/remote"
	"github.com/SimonXming/circle/store"
	"github.com/SimonXming/circle/utils"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func NewEchoServer() *echo.Echo {
	return echo.New()
}

func GetRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func PostScmAccount(c echo.Context) error {
	in := new(model.ScmAccount)
	if err := c.Bind(in); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	account := &model.ScmAccount{
		Host:         in.Host,
		Login:        in.Login,
		Password:     in.Password,
		Type:         in.Type,
		PrivateToken: in.PrivateToken,
	}

	err := store.ScmAccountCreate(c, account)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	return nil

}

func GetScmAccounts(c echo.Context) error {
	accounts, err := store.ScmAccountList(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	return c.JSON(http.StatusOK, accounts)

}

func GetScmAccount(c echo.Context) error {
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
	return c.JSON(http.StatusOK, account)

}

func GetRemoteRepos(c echo.Context) error {
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
	err = utils.SetupRemote(c, account)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	remote := remote.FromContext(c)
	repos, err := remote.Repos()
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

	account, err := store.ScmAccountLoad(c, scmId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	err = utils.SetupRemote(c, account)

	r := &model.Repo{
		ScmId:  scmId,
		Owner:  in.Owner,
		Name:   in.Name,
		Clone:  in.Clone,
		Branch: in.Branch,
	}

	link := "http://localhost:8080/some/hook/url?access_token=123"

	// activate the repository before we make any
	// local changes to the database.
	err = remote.Activate(c, r, link)
	if err != nil {
		c.String(500, err.Error())
		return err
	}

	err = store.RepoCreate(c, r)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	return nil
}

func PostSecret(c echo.Context) error {
	in := new(model.Secret)
	if err := c.Bind(in); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	secret := &model.Secret{
		RepoID: 123,
		Name:   in.Name,
		Value:  in.Value,
	}

	err := store.SecretCreate(c, secret)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	return nil
}
