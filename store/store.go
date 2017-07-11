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
	ConfigFind(*model.Repo) (*model.Config, error)

	RepoCreate(*model.Repo) error
	RepoLoad(int64) (*model.Repo, error)

	// GetRepoName gets a repo by its full name.
	GetRepoName(string) (*model.Repo, error)
	GetRepoScmName(int64, string) (*model.Repo, error)

	SecretCreate(*model.Secret) error

	BuildCreate(*model.Build, ...*model.Proc) error

	ProcCreate([]*model.Proc) error
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

func ConfigCreate(c echo.Context, conf *model.Config) error {
	return FromContext(c).ConfigCreate(conf)
}

func ConfigFind(c echo.Context, repo *model.Repo) (*model.Config, error) {
	return FromContext(c).ConfigFind(repo)
}

func RepoCreate(c echo.Context, repo *model.Repo) error {
	return FromContext(c).RepoCreate(repo)
}

func RepoLoad(c echo.Context, id int64) (*model.Repo, error) {
	return FromContext(c).RepoLoad(id)
}

func GetRepoOwnerName(c echo.Context, owner, name string) (*model.Repo, error) {
	return FromContext(c).GetRepoName(owner + "/" + name)
}

func GetRepoScmIDOwnerName(c echo.Context, scmID int64, owner, name string) (*model.Repo, error) {
	return FromContext(c).GetRepoScmName(scmID, owner+"/"+name)
}

func SecretCreate(c echo.Context, secret *model.Secret) error {
	return FromContext(c).SecretCreate(secret)
}

func BuildCreate(c echo.Context, build *model.Build, procs ...*model.Proc) error {
	return FromContext(c).BuildCreate(build, procs...)
}

func ProcCreate(c echo.Context, procs []*model.Proc) error {
	return FromContext(c).ProcCreate(procs)
}
