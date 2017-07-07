package handler

import (
	// "fmt"
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/store"
	"github.com/labstack/echo"
	"net/http"
)

func NewEchoServer() *echo.Echo {
	return echo.New()
}

func GetRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func PostRepo(c echo.Context) error {
	in := new(model.Repo)
	if err := c.Bind(in); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	repo := &model.Repo{
		Clone:  in.Clone,
		Kind:   in.Kind,
		Branch: in.Branch,
	}

	err := store.RepoCreate(c, repo)
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
