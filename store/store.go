package store

import (
	// "io"

	"github.com/SimonXming/circle/model"
	"github.com/labstack/echo"
)

type Store interface {
	ScmAccountCreate(*model.ScmAccount) error
	ScmAccountList() ([]*model.ScmAccount, error)
	ScmAccountLoad(int64) (*model.ScmAccount, error)

	ConfigCreate(*model.Config) error
	ConfigLoad(int64) (*model.Config, error)

	RepoCreate(*model.Repo) error

	// GetRepoName gets a repo by its full name.
	GetRepoName(string) (*model.Repo, error)

	SecretCreate(*model.Secret) error

	BuildCreate(*model.Build) error
}

func ScmAccountCreate(c echo.Context, account *model.ScmAccount) error {
	return FromContext(c).ScmAccountCreate(account)
}

func ScmAccountList(c echo.Context) ([]*model.ScmAccount, error) {
	return FromContext(c).ScmAccountList()
}

func ScmAccountLoad(c echo.Context, id int64) (*model.ScmAccount, error) {
	return FromContext(c).ScmAccountLoad(id)
}

func RepoCreate(c echo.Context, repo *model.Repo) error {
	return FromContext(c).RepoCreate(repo)
}

func GetRepoOwnerName(c echo.Context, owner, name string) (*model.Repo, error) {
	return FromContext(c).GetRepoName(owner + "/" + name)
}

func SecretCreate(c echo.Context, secret *model.Secret) error {
	return FromContext(c).SecretCreate(secret)
}
