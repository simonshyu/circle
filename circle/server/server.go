package server

import (
	"github.com/labstack/echo"
	"github.com/urfave/cli"
	"golang.org/x/net/websocket"
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
	e.WebSocket("/ws", 
	e.Logger.Fatal(e.Start(":1323"))
	return nil
}

func getRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func Broker(c echo.Context) error {
	ws := c.Socket()
	msg := ""

	for {
		if err = websocket.Message.Send(ws, "Hello, Client!"); err != nil {
			return
		}
		if err = websocket.Message.Receive(ws, &msg); err != nil {
			return
		}
		fmt.Println(msg)
	}
}