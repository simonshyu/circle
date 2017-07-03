package router

import (
	"github.com/SimonXming/circle/handler"
	"github.com/labstack/echo"
)

func Load(e *echo.Echo) {
	e.GET("/", handler.GetRoot)
	e.GET("/ws/broker", handler.Broker)
}
