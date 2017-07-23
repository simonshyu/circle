package gitlab

import (
	"bytes"
	"net/http"
	"testing"

	// "github.com/SimonXming/circle/model"
	"github.com/SimonXming/circle/remote/gitlab/testdata"
	"github.com/franela/goblin"
)

func Test_Gitlab(t *testing.T) {
	// setup a dummy github server
	var server = testdata.NewServer()
	defer server.Close()

	env := server.URL + "?client_id=test&client_secret=test"

	gitlab := Load(env)

	// var repo = model.Repo{
	// 	Name:  "diaspora-client",
	// 	Owner: "diaspora",
	// }
	g := goblin.Goblin(t)
	g.Describe("Gitlab Plugin", func() {
		g.Describe("Hook", func() {
			g.Describe("Push hook", func() {
				g.It("Should parse actual push hoook", func() {
					req, _ := http.NewRequest(
						"POST",
						"http://localhost:8000/api/hook?scm_id=1&owner=diaspora&name=diaspora-client",
						bytes.NewReader(testdata.PushHook),
					)

					repo, build, err := gitlab.Hook(req)

					if err != nil {
						println(err.Error())
					}
					// g.Assert(err == nil).IsTrue()
					g.Assert(repo.Owner).Equal("mike")
					g.Assert(repo.Name).Equal("diaspora")
					g.Assert(repo.Branch).Equal("develop")
					g.Assert(build.Ref).Equal("refs/heads/master")
				})
			})
			g.Describe("Tag push hook", func() {
				g.It("Should parse tag push hook", func() {
					req, _ := http.NewRequest(
						"POST",
						"http://localhost:8000/api/hook?scm_id=1&owner=diaspora&name=diaspora-client",
						bytes.NewReader(testdata.TagHook),
					)

					repo, build, err := gitlab.Hook(req)
					if err != nil {
						println(err.Error())
					}

					// g.Assert(err == nil).IsTrue()
					g.Assert(repo.Owner).Equal("jsmith")
					g.Assert(repo.Name).Equal("example")
					// g.Assert(repo.Avatar).Equal("http://example.com/uploads/project/avatar/555/Outh-20-Logo.jpg")
					g.Assert(repo.Branch).Equal("develop")
					g.Assert(build.Ref).Equal("refs/tags/v1.0.0")

				})
			})
			g.Describe("Merge request hook", func() {
				g.It("Should parse merge request hook", func() {
					req, _ := http.NewRequest(
						"POST",
						"http://localhost:8000/api/hook?scm_id=1&owner=diaspora&name=diaspora-client",
						bytes.NewReader(testdata.MergeRequestHook),
					)

					repo, build, err := gitlab.Hook(req)
					if err != nil {
						println(err.Error())
					}
					// g.Assert(err == nil).IsTrue()
					// g.Assert(repo.Avatar).Equal("http://example.com/uploads/project/avatar/555/Outh-20-Logo.jpg")
					g.Assert(repo.Branch).Equal("develop")
					g.Assert(repo.Owner).Equal("awesome_space")
					g.Assert(repo.Name).Equal("awesome_project")

					g.Assert(build.Branch).Equal("ms-viewport")
				})
			})
		})
	})
}
