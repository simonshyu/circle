package server

import (
	// "fmt"
	ciserver "github.com/SimonXming/circle/handler"
	"github.com/SimonXming/circle/router"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:   "server",
	Usage:  "starts the circle server daemon",
	Action: server,
}

func server(c *cli.Context) error {
	e := ciserver.NewEchoServer()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	router.Load(e)
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
	return nil
}
