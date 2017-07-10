package remote

//go:generate mockery -name Remote -output mock -case=underscore

import (
	"net/http"
	// "time"
	"github.com/labstack/echo"

	"github.com/SimonXming/circle/model"
)

type Remote interface {
	// Repo fetches the named repository from the remote system.
	Repo(owner, repo string) (*model.Repo, error)

	// Repos fetches a list of repos from the remote system.
	Repos() ([]*model.Repo, error)

	// private repositories from a remote system.
	Netrc(r *model.Repo) (*model.Netrc, error)

	// Activate activates a repository by creating the post-commit hook.
	Activate(r *model.Repo, link string) error

	// Deactivate deactivates a repository by removing all previously created
	// post-commit hooks matching the given link.
	Deactivate(r *model.Repo, link string) error

	// Hook parses the post-commit hook from the Request body and returns the
	// required data in a standard format.
	Hook(r *http.Request) (*model.Repo, *model.Build, error)
}

// Repo fetches the named repository from the remote system.
func Repos(c echo.Context) ([]*model.Repo, error) {
	return FromContext(c).Repos()
}

// Repo fetches the named repository from the remote system.
func Activate(c echo.Context, r *model.Repo, link string) error {
	return FromContext(c).Activate(r, link)
}
