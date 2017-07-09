package utils

import (
	"fmt"
	"github.com/labstack/echo"

	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/remote"
	"github.com/SimonXming/circle/remote/gitlab"
)

func SetupRemote(c echo.Context, account *model.ScmAccount) error {
	r, err := FromScmAccount(account)
	if err != nil {
		return err
	}
	remote.ToContext(c, r)
	return nil
}

func FromScmAccount(account *model.ScmAccount) (remote.Remote, error) {
	switch account.Type {
	case "gitlab":
		return setupGitlab(account)
	default:
		return nil, fmt.Errorf("version control system not configured")
	}
}

// helper function to setup the Gitlab remote from the CLI arguments.
func setupGitlab(account *model.ScmAccount) (remote.Remote, error) {
	println("account", account.PrivateToken)
	return gitlab.New(gitlab.Opts{
		URL:          account.Host,
		Username:     account.Login,
		Password:     account.Password,
		PrivateToken: account.PrivateToken,
		SkipVerify:   false,
	})
}
