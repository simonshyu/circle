package router

import (
	"github.com/SimonXming/circle/handler"
	"github.com/labstack/echo"
)

func Load(e *echo.Echo) {
	e.GET("/", handler.GetRoot)
	// e.POST("/repo", handler.PostRepo)
	e.POST("/scm_account", handler.PostScmAccount)
	e.GET("/ws/broker", handler.RPCHandler)
}
