package router

import (
	"github.com/SimonXming/circle/handler"
	"github.com/labstack/echo"
)

func Load(e *echo.Echo) {
	e.GET("/", handler.GetRoot)
	e.GET("/queue/info", handler.GetQueueInfo)

	scmGroup := e.Group("/scm")
	{
		scmGroup.POST("", handler.PostScmAccount)
		scmGroup.GET("", handler.GetScmAccounts)
		scmGroup.GET("/repo", handler.GetAllRepo)
		scmGroup.GET("/:scmID", handler.GetScmAccount)
		scmGroup.GET("/:scmID/remote", handler.GetRemoteRepos)

		repoGroup := scmGroup.Group("/:scmID/repo")
		{
			repoGroup.POST("", handler.PostRepo)
			repoGroup.GET("", handler.GetRepos)
			repoGroup.POST("/:repoID/config", handler.PostConfig)
			repoGroup.POST("/:repoID/build", handler.PostBuild)
		}
	}

	websocketGroup := e.Group("/ws")
	websocketGroup.GET("/broker", handler.RPCHandler)

	e.POST("/hook", handler.PostHook)
	e.POST("/api/hook", handler.PostHook)
}
