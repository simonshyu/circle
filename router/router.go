package router

import (
	"github.com/SimonXming/circle/handler"
	"github.com/labstack/echo"
)

func Load(e *echo.Echo) {
	e.GET("/", handler.GetRoot)
	// e.POST("/repo", handler.PostRepo)
	e.POST("/scm", handler.PostScmAccount)
	e.GET("/scm", handler.GetScmAccounts)
	e.GET("/scm/:scmID", handler.GetScmAccount)
	e.GET("/scm/:scmID/repos/remote", handler.GetRemoteRepos)
	e.POST("/scm/:scmID/repos/:owner/:name", handler.PostRepo)
	e.POST("/scm/:scmID/repos/:repoID/config", handler.PostConfig)
	e.GET("/ws/broker", handler.RPCHandler)
}
