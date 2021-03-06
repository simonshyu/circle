package store

import (
	"io"

	"github.com/labstack/echo"
	"github.com/simonshyu/circle/model"
	"github.com/simonshyu/circle/utils"
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
	RepoList() ([]*model.Repo, error)
	RepoFind(*model.ScmAccount) ([]*model.Repo, error)
	GetRepoScmName(int64, string) (*model.Repo, error)

	SecretCreate(*model.Secret) error
	SecretUpdate(*model.Secret) error

	BuildCreate(*model.Build, ...*model.Proc) error
	BuildFind(*model.Repo) ([]*model.Build, error)
	BuildLoad(int64) (*model.Build, error)
	BuildUpdate(*model.Build) error
	GetBuildNumber(*model.Repo, int) (*model.Build, error)

	ProcCreate([]*model.Proc) error
	ProcList(*model.Build) ([]*model.Proc, error)
	ProcLoad(int64) (*model.Proc, error)
	ProcChild(*model.Build, int, string) (*model.Proc, error)
	ProcUpdate(*model.Proc) error
	ProcClear(*model.Build) error

	TaskList() ([]*model.Task, error)
	TaskInsert(*model.Task) error
	TaskDelete(string) error

	LogFind(*model.Proc) (io.ReadCloser, error)
	LogSave(*model.Proc, io.Reader) error

	RegistryCreate(*model.Registry) error
	RegistryUpdate(*model.Registry) error
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

func RepoList(c echo.Context) ([]*model.Repo, error) {
	return FromContext(c).RepoList()
}

func RepoFind(c echo.Context, scm *model.ScmAccount) ([]*model.Repo, error) {
	return FromContext(c).RepoFind(scm)
}

func RepoLoad(c echo.Context, id int64) (*model.Repo, error) {
	return FromContext(c).RepoLoad(id)
}

func GetRepoScmIDOwnerName(c echo.Context, scmID int64, owner, name string) (*model.Repo, error) {
	return FromContext(c).GetRepoScmName(scmID, owner+"/"+name)
}

func SecretCreate(c echo.Context, secret *model.Secret) error {
	return FromContext(c).SecretCreate(secret)
}

func SecretUpdate(c echo.Context, secret *model.Secret) error {
	return FromContext(c).SecretUpdate(secret)
}

func BuildCreate(c echo.Context, build *model.Build, procs ...*model.Proc) error {
	return FromContext(c).BuildCreate(build, procs...)
}

func BuildFind(c echo.Context, repo *model.Repo) ([]*model.Build, error) {
	return FromContext(c).BuildFind(repo)
}

func BuildLoad(c echo.Context, id int64) (*model.Build, error) {
	return FromContext(c).BuildLoad(id)
}

func BuildUpdate(c echo.Context, build *model.Build) error {
	return FromContext(c).BuildUpdate(build)
}

func GetBuildNumber(c echo.Context, repo *model.Repo, num int) (*model.Build, error) {
	return FromContext(c).GetBuildNumber(repo, num)
}

func ProcCreate(c echo.Context, procs []*model.Proc) error {
	return FromContext(c).ProcCreate(procs)
}

func ProcList(c echo.Context, build *model.Build) ([]*model.Proc, error) {
	return FromContext(c).ProcList(build)
}

func ProcLoad(c echo.Context, id int64) (*model.Proc, error) {
	return FromContext(c).ProcLoad(id)
}

func ProcChild(c echo.Context, build *model.Build, pid int, child string) (*model.Proc, error) {
	return FromContext(c).ProcChild(build, pid, child)
}

func ProcUpdate(c echo.Context, proc *model.Proc) error {
	return FromContext(c).ProcUpdate(proc)
}

func ProcClear(c echo.Context, build *model.Build) error {
	return FromContext(c).ProcClear(build)
}

func LogFind(c echo.Context, proc *model.Proc) (io.ReadCloser, error) {
	return FromContext(c).LogFind(proc)
}

func RegistryCreate(c echo.Context, registry *model.Registry) error {
	return FromContext(c).RegistryUpdate(registry)
}

func RegistryUpdate(c echo.Context, registry *model.Registry) error {
	return FromContext(c).RegistryUpdate(registry)
}

// helper: 合并 ScmAccountLoad 和 SetupRemote 的功能
func SetupRemoteWithScmID(c echo.Context, id int64) (*model.ScmAccount, error) {
	account, err := FromContext(c).ScmAccountLoad(id)
	if err != nil {
		return nil, err
	}
	err = utils.SetupRemote(c, account)
	return account, err
}
