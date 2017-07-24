package gitlab

import (
	"fmt"
	"github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/remote"
	"github.com/SimonXming/circle/remote/gitlab/client"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
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

func Load(config string) *Gitlab {
	url_, err := url.Parse(config)
	if err != nil {
		panic(err)
	}
	params := url_.Query()
	url_.RawQuery = ""

	gitlab := Gitlab{}
	gitlab.URL = url_.String()
	gitlab.Machine = "172.24.6.123"
	gitlab.Username = params.Get("username")
	gitlab.Password = params.Get("password")
	// gitlab.AllowedOrgs = params["orgs"]
	gitlab.SkipVerify, _ = strconv.ParseBool(params.Get("skip_verify"))
	gitlab.HideArchives, _ = strconv.ParseBool(params.Get("hide_archives"))
	// gitlab.Open, _ = strconv.ParseBool(params.Get("open"))

	// switch params.Get("clone_mode") {
	// case "oauth":
	// 	gitlab.CloneMode = "oauth"
	// default:
	// 	gitlab.CloneMode = "token"
	// }

	// this is a temp workaround
	// gitlab.Search, _ = strconv.ParseBool(params.Get("search"))

	return &gitlab
}

// Repo fetches the named repository from the remote system.
func (g *Gitlab) Repo(owner, name string) (*model.Repo, error) {
	client := NewClient(g.URL, g.PrivateToken, g.SkipVerify)
	id, err := GetProjectId(g, client, owner, name)
	if err != nil {
		return nil, err
	}
	repo_, err := client.Project(id)
	if err != nil {
		return nil, err
	}

	repo := &model.Repo{}
	repo.Owner = owner
	repo.Name = name
	repo.FullName = repo_.PathWithNamespace
	repo.Link = repo_.Url
	repo.Clone = repo_.HttpRepoUrl
	repo.Branch = "master"

	if repo_.DefaultBranch != "" {
		repo.Branch = repo_.DefaultBranch
	}

	repo.IsPrivate = !repo_.Public

	// if g.PrivateMode {
	// 	repo.IsPrivate = true
	// } else {
	// 	repo.IsPrivate = !repo_.Public
	// }

	return repo, err
}

func (g *Gitlab) Repos() ([]*model.RepoLite, error) {
	client := NewClient(g.URL, g.PrivateToken, g.SkipVerify)

	var repos = []*model.RepoLite{}

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

		repos = append(repos, &model.RepoLite{
			Owner:    owner,
			Name:     name,
			FullName: repo.PathWithNamespace,
			Clone:    clone_url,
			Branch:   branch,
		})
	}

	return repos, err
}

func (g *Gitlab) Netrc(r *model.ScmAccount) (*model.Netrc, error) {
	gitlabUrl, err := url.Parse(r.Host)
	if err != nil {
		return nil, err
	}
	return &model.Netrc{
		Login:    r.Login,
		Password: r.Password,
		Machine:  gitlabUrl.Hostname(),
	}, nil
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

	hookScmID := uri.Query().Get("scm_id")
	hookUrl := fmt.Sprintf("%s://%s%s?scmID=%s", uri.Scheme, uri.Host, uri.Path, hookScmID)
	hookToken := uri.Query().Get("access_token")
	ssl_verify := g.SkipVerify
	all_push := r.AllowPush

	err = client.AddProjectHook(id, hookUrl, hookToken, all_push, ssl_verify)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gitlab) Deactivate(r *model.Repo, link string) error {
	return nil
}

// ParseHook parses the post-commit hook from the Request body
// and returns the required data in a standard format.
func (g *Gitlab) Hook(req *http.Request) (*model.Repo, *model.Build, error) {
	defer req.Body.Close()
	var payload, _ = ioutil.ReadAll(req.Body)
	var parsed, err = client.ParseHook(payload)
	if err != nil {
		return nil, nil, err
	}

	switch parsed.ObjectKind {
	case "merge_request":
		return mergeRequest(parsed, req)
		return nil, nil, nil
	case "tag_push", "push":
		return push(parsed, req)
	default:
		return nil, nil, nil
	}

}

func mergeRequest(parsed *client.HookPayload, req *http.Request) (*model.Repo, *model.Build, error) {

	repo := &model.Repo{}

	obj := parsed.ObjectAttributes
	if obj == nil {
		return nil, nil, fmt.Errorf("object_attributes key expected in merge request hook")
	}

	target := obj.Target
	source := obj.Source

	if target == nil && source == nil {
		return nil, nil, fmt.Errorf("target and source keys expected in merge request hook")
	} else if target == nil {
		return nil, nil, fmt.Errorf("target key expected in merge request hook")
	} else if source == nil {
		return nil, nil, fmt.Errorf("source key exptected in merge request hook")
	}

	if target.PathWithNamespace != "" {
		var err error
		if repo.Owner, repo.Name, err = ExtractFromPath(target.PathWithNamespace); err != nil {
			return nil, nil, err
		}
		repo.FullName = target.PathWithNamespace
	} else {
		repo.Owner = req.FormValue("owner")
		repo.Name = req.FormValue("name")
		repo.FullName = fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
	}

	repo.Link = target.WebUrl

	if target.GitHttpUrl != "" {
		repo.Clone = target.GitHttpUrl
	} else {
		repo.Clone = target.HttpUrl
	}

	if target.DefaultBranch != "" {
		repo.Branch = target.DefaultBranch
	} else {
		repo.Branch = "master"
	}

	// if target.AvatarUrl != "" {
	// 	repo.Avatar = target.AvatarUrl
	// }

	build := &model.Build{}
	build.Event = "pull_request"

	lastCommit := obj.LastCommit
	if lastCommit == nil {
		return nil, nil, fmt.Errorf("last_commit key expected in merge request hook")
	}

	// build.Message = lastCommit.Message
	build.Commit = lastCommit.Id
	build.Remote = source.HttpUrl

	build.Ref = fmt.Sprintf("refs/merge-requests/%d/head", obj.IId)

	build.Branch = obj.SourceBranch

	author := lastCommit.Author
	if author == nil {
		return nil, nil, fmt.Errorf("author key expected in merge request hook")
	}

	// build.Author = author.Name
	// build.Email = author.Email

	// if len(build.Email) != 0 {
	// 	build.Avatar = GetUserAvatar(build.Email)
	// }

	// build	.Title = obj.Title
	build.Link = obj.Url

	return repo, build, nil
}

func push(parsed *client.HookPayload, req *http.Request) (*model.Repo, *model.Build, error) {
	repo := &model.Repo{}

	// Since gitlab 8.5, used project instead repository key
	// see https://gitlab.com/gitlab-org/gitlab-ce/blob/master/doc/web_hooks/web_hooks.md#web-hooks
	if project := parsed.Project; project != nil {
		var err error
		if repo.Owner, repo.Name, err = ExtractFromPath(project.PathWithNamespace); err != nil {
			return nil, nil, err
		}

		// repo.Avatar = project.AvatarUrl
		repo.Link = project.WebUrl
		repo.Clone = project.GitHttpUrl
		repo.FullName = project.PathWithNamespace
		repo.Branch = project.DefaultBranch

		switch project.VisibilityLevel {
		case 0:
			repo.IsPrivate = true
		case 10:
			repo.IsPrivate = true
		case 20:
			repo.IsPrivate = false
		}
	} else if repository := parsed.Repository; repository != nil {
		repo.Owner = req.FormValue("owner")
		repo.Name = req.FormValue("name")
		repo.Link = repository.URL
		repo.Clone = repository.GitHttpUrl
		repo.Branch = "master"
		repo.FullName = fmt.Sprintf("%s/%s", req.FormValue("owner"), req.FormValue("name"))

		switch repository.VisibilityLevel {
		case 0:
			repo.IsPrivate = true
		case 10:
			repo.IsPrivate = true
		case 20:
			repo.IsPrivate = false
		}
	} else {
		return nil, nil, fmt.Errorf("No project/repository keys given")
	}

	build := &model.Build{}
	build.Event = model.EventPush
	build.Commit = parsed.After
	build.Branch = parsed.Branch()
	build.Ref = parsed.Ref
	// hook.Commit.Remote = cloneUrl

	// var head = parsed.Head()
	// build.Message = head.Message
	// build.Timestamp = head.Timestamp

	// extracts the commit author (ideally email)
	// from the post-commit hook
	// switch {
	// case head.Author != nil:
	// 	build.Email = head.Author.Email
	// 	build.Author = parsed.UserName
	// 	if len(build.Email) != 0 {
	// 		build.Avatar = GetUserAvatar(build.Email)
	// 	}
	// case head.Author == nil:
	// 	build.Author = parsed.UserName
	// }

	if strings.HasPrefix(build.Ref, "refs/tags/") {
		build.Event = model.EventTag
	}

	return repo, build, nil
}
