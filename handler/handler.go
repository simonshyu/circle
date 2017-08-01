package handler

import (
	"crypto/sha256"
	"fmt"

	"context"
	"github.com/simonshyu/circle/model"
	"github.com/simonshyu/circle/store"
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

func GetQueueInfo(c echo.Context) error {
	info := Config.Services.Queue.Info(context.Background())
	return c.JSON(http.StatusOK, info.Stats)
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

func GetConfig(c echo.Context) error {
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
	config, err := store.ConfigFind(c, repo)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	return c.JSON(http.StatusOK, config)

}

func PostConfig(c echo.Context) error {
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

	in := new(model.Config)
	if err := c.Bind(in); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	config := &model.Config{
		RepoID: repo.ID,
		Data:   in.Data,
	}
	config.Hash = shasum([]byte(in.Data))

	err = store.ConfigCreate(c, config)
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

func shasum(raw []byte) string {
	sum := sha256.Sum256(raw)
	return fmt.Sprintf("%x", sum)
}
