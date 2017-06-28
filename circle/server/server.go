package server

import (
	"github.com/labstack/echo"
	"github.com/urfave/cli"
	"net/http"
)

var Command = cli.Command{
	Name:   "server",
	Usage:  "starts the circle server daemon",
	Action: server,
}

func server(c *cli.Context) error {
	e := echo.New()
	e.GET("/", getRoot)
	e.Logger.Fatal(e.Start(":1323"))
	return nil
}

func getRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
