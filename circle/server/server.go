package server

import (
	// "fmt"
	ciserver "github.com/SimonXming/circle/handler"
	"github.com/SimonXming/circle/router"
	cimiddleware "github.com/SimonXming/circle/router/middleware"
	"github.com/SimonXming/circle/store"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
)

var Command = cli.Command{
	Name:   "server",
	Usage:  "starts the circle server daemon",
	Action: server,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "driver",
			Usage: "database driver",
			Value: "mysql",
		},
		cli.StringFlag{
			Name:  "datasource",
			Usage: "database driver configuration string",
			Value: "test:test@tcp(127.0.0.1:3306)/test",
		},
	},
}

func server(c *cli.Context) error {
	s := setupStore(c)
	setupEvilGlobals(c, s)

	e := ciserver.NewEchoServer()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(cimiddleware.StoreWithConfig(s))
	router.Load(e)
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
	return nil
}

func setupEvilGlobals(c *cli.Context, v store.Store) {
	// storage
	ciserver.Config.Storage.Config = v
}
