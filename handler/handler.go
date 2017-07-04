package handler

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/sourcegraph/jsonrpc2"
	websocketrpc "github.com/sourcegraph/jsonrpc2/websocket"
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

func NewEchoServer() *echo.Echo {
	return echo.New()
}

func GetRoot(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func Broker(c echo.Context) error {
	// receive agent connection
	ws, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	// get jsonrpc2 websokcet connection
	conn := jsonrpc2.NewConn(ctx,
		websocketrpc.NewObjectStream(ws),
		jsonrpc2.HandlerWithError(router),
	)
	defer func() {
		cancel()
		conn.Close()
	}()
	<-conn.DisconnectNotify()
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
