package gitlab

import (
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/remote"
	"net"
	"net/http"
	"net/url"
)

// Opts defines configuration options.
type Opts struct {
	URL        string // Gogs server url.
	Username   string // Optional machine account username.
	Password   string // Optional machine account password.
	SkipVerify bool   // Skip ssl verification.
}

// New returns a Remote implementation that integrates with Gitlab, an open
// source Git service. See https://gitlab.com
func New(opts Opts) (remote.Remote, error) {
	url, err := url.Parse(opts.URL)
	if err != nil {
		return nil, err
	}
	host, _, err := net.SplitHostPort(url.Host)
	if err == nil {
		url.Host = host
	}
	return &Gitlab{
		URL:        opts.URL,
		Machine:    url.Host,
		Username:   opts.Username,
		Password:   opts.Password,
		SkipVerify: opts.SkipVerify,
	}, nil
}

type Gitlab struct {
	URL          string
	Machine      string
	Username     string
	Password     string
	SkipVerify   bool
	HideArchives bool
	Search       bool
}

func (g *Gitlab) Repo(owner, repo string) (*model.Repo, error) {
	return nil, nil
}

func (g *Gitlab) Repos() ([]*model.Repo, error) {
	return nil, nil
}

func (g *Gitlab) Netrc(r *model.Repo) (*model.Netrc, error) {
	return nil, nil
}

func (g *Gitlab) Activate(r *model.Repo, link string) error {
	return nil
}

func (g *Gitlab) Deactivate(r *model.Repo, link string) error {
	return nil
}

func (g *Gitlab) Hook(r *http.Request) (*model.Repo, *model.Build, error) {
	return nil, nil, nil
}
