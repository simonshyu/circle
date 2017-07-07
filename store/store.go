package store

import (
	// "io"

	"github.com/SimonXming/circle/model"
	"github.com/labstack/echo"
)

type Store interface {
	ScmAccountCreate(*model.ScmAcount) error

	ConfigCreate(*model.Config) error
	ConfigLoad(int64) (*model.Config, error)

	RepoCreate(*model.Repo) error

	SecretCreate(*model.Secret) error

	BuildCreate(*model.Build) error
}

func ScmAccountCreate(c echo.Context, account *model.ScmAcount) error {
	return FromContext(c).ScmAccountCreate(account)
}

func RepoCreate(c echo.Context, repo *model.Repo) error {
	return FromContext(c).RepoCreate(repo)
}

func SecretCreate(c echo.Context, secret *model.Secret) error {
	return FromContext(c).SecretCreate(secret)
}
