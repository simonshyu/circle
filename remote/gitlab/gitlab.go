package gitlab

import (
	"fmt"
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/remote"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// Opts defines configuration options.
type Opts struct {
	URL          string // Gogs server url.
	PrivateToken string // Gitlab PrivateToken
	Username     string // Optional machine account username.
	Password     string // Optional machine account password.
	SkipVerify   bool   // Skip ssl verification.
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
	fmt.Printf("%v", opts)
	return &Gitlab{
		URL:          opts.URL,
		Machine:      url.Host,
		Username:     opts.Username,
		Password:     opts.Password,
		PrivateToken: opts.PrivateToken,
		SkipVerify:   opts.SkipVerify,
	}, nil
}

type Gitlab struct {
	URL          string
	Machine      string
	Username     string
	Password     string
	PrivateToken string
	SkipVerify   bool
	HideArchives bool
	Search       bool
}

func (g *Gitlab) Repo(owner, repo string) (*model.Repo, error) {
	return nil, nil
}

func (g *Gitlab) Repos() ([]*model.Repo, error) {
	client := NewClient(g.URL, g.PrivateToken, g.SkipVerify)

	var repos = []*model.Repo{}

	all, err := client.AllProjects(g.HideArchives)
	if err != nil {
		return repos, err
	}

	for _, repo := range all {
		var parts = strings.Split(repo.PathWithNamespace, "/")
		var owner = parts[0]
		var name = parts[1]
		var clone_url = repo.HttpRepoUrl
		var branch = repo.DefaultBranch

		repos = append(repos, &model.Repo{
			Owner:  owner,
			Name:   name,
			Clone:  clone_url,
			Branch: branch,
		})
	}

	return repos, err
}

func (g *Gitlab) Netrc(r *model.Repo) (*model.Netrc, error) {
	return nil, nil
}

func (g *Gitlab) Activate(r *model.Repo, link string) error {
	client := NewClient(g.URL, g.PrivateToken, g.SkipVerify)
	id, err := GetProjectId(g, client, r.Owner, r.Name)
	if err != nil {
		return err
	}
	uri, err := url.Parse(link)
	if err != nil {
		return err
	}

	// circleUrl := "http://localhost:8000/hook?access_token=1234"
	circleUrl := fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)
	circleToken := uri.Query().Get("access_token")
	ssl_verify := g.SkipVerify
	all_push := r.AllowPush

	err = client.AddProjectHook(id, circleUrl, circleToken, all_push, ssl_verify)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gitlab) Deactivate(r *model.Repo, link string) error {
	return nil
}

func (g *Gitlab) Hook(r *http.Request) (*model.Repo, *model.Build, error) {
	return nil, nil, nil
}
