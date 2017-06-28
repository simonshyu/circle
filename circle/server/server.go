package server

import (
	"context"
	"errors"
	// "fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sourcegraph/jsonrpc2"
	websocketrpc "github.com/sourcegraph/jsonrpc2/websocket"
	"github.com/urfave/cli"
	"net/http"
)

const (
	methodNext   = "next"
	methodWait   = "wait"
	methodInit   = "init"
	methodDone   = "done"
	methodExtend = "extend"
	methodUpdate = "update"
	methodUpload = "upload"
	methodLog    = "log"
)

var (
	errNoSuchMethod = errors.New("No such rpc method")
	upgrader        = websocket.Upgrader{}
)

var Command = cli.Command{
	Name:   "server",
	Usage:  "starts the circle server daemon",
	Action: server,
}

func server(c *cli.Context) error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", getRoot)
	e.GET("/ws/broker", broker)
	e.Logger.Fatal(e.Start("0.0.0.0:8000"))
	return nil
}

func getRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func broker(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	conn := jsonrpc2.NewConn(ctx,
		websocketrpc.NewObjectStream(ws),
		jsonrpc2.HandlerWithError(router),
	)
	defer func() {
		cancel()
		conn.Close()
	}()
	<-conn.DisconnectNotify()
	// defer ws.Close()

	// for {
	// 	// Write
	// 	err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
	// 	if err != nil {
	// 		c.Logger().Error(err)
	// 	}

	// 	// Read
	// 	_, msg, err := ws.ReadMessage()
	// 	if err != nil {
	// 		c.Logger().Error(err)
	// 	}
	// 	fmt.Printf("%s\n", msg)
	// }
	return nil
}

func router(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	switch req.Method {
	case methodNext:
		println("methodNext")
		return nil, errNoSuchMethod
		// return s.next(ctx, req)
	case methodWait:
		println("methodWait")
		return nil, errNoSuchMethod
		// return s.wait(ctx, req)
	case methodInit:
		println("methodInit")
		return nil, errNoSuchMethod
		// return s.init(ctx, req)
	case methodDone:
		println("methodDone")
		return nil, errNoSuchMethod
		// return s.done(ctx, req)
	case methodExtend:
		println("methodExtend")
		return nil, errNoSuchMethod
		// return s.extend(ctx, req)
	case methodUpdate:
		println("methodExtend")
		return nil, errNoSuchMethod
		// return s.update(req)
	case methodLog:
		println("methodLog")
		return nil, errNoSuchMethod
		// return s.log(req)
	case methodUpload:
		println("methodUpload")
		return nil, errNoSuchMethod
		// return s.upload(req)
	default:
		return nil, errNoSuchMethod
	}
}
