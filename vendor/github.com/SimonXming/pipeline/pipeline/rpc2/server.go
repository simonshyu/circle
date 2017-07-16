package rpc2

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	websocketrpc "github.com/sourcegraph/jsonrpc2/websocket"
)

// errNoSuchMethod is returned when the name rpc method does not exist.
var errNoSuchMethod = errors.New("No such rpc method")

// noContext is an empty context used when no context is required.
var noContext = context.Background()

// Server represents an rpc server.
type Server struct {
	peer Peer
}

// NewServer returns an rpc Server.
func NewServer(peer Peer) *Server {
	return &Server{peer}
}

// ServeHTTP implements an http.Handler that answers rpc requests.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	conn := jsonrpc2.NewConn(ctx,
		websocketrpc.NewObjectStream(c),
		jsonrpc2.HandlerWithError(s.router),
	)
	defer func() {
		cancel()
		conn.Close()
	}()
	<-conn.DisconnectNotify()
}

// router implements an jsonrpc2.Handler that answers RPC requests.
func (s *Server) router(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	switch req.Method {
	case methodNext:
		return s.next(ctx, req)
	case methodWait:
		return s.wait(ctx, req)
	case methodInit:
		return s.init(ctx, req)
	case methodDone:
		return s.done(ctx, req)
	case methodExtend:
		return s.extend(ctx, req)
	default:
		return nil, errNoSuchMethod
	}
}

// next unmarshals the rpc request parameters and invokes the peer.Next
// procedure. The results are retuned and written to the rpc response.
func (s *Server) next(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	in := Filter{}
	if err := json.Unmarshal([]byte(*req.Params), &in); err != nil {
		return nil, err
	}
	return s.peer.Next(ctx, in)
}

// wait unmarshals the rpc request parameters and invokes the peer.Wait
// procedure. The results are retuned and written to the rpc response.
func (s *Server) wait(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	var id string
	err := json.Unmarshal([]byte(*req.Params), &id)
	if err != nil {
		return nil, err
	}
	return nil, s.peer.Wait(ctx, id)
}

// init unmarshals the rpc request parameters and invokes the peer.Init
// procedure. The results are retuned and written to the rpc response.
func (s *Server) init(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	in := new(updateReq)
	if err := json.Unmarshal([]byte(*req.Params), in); err != nil {
		return nil, err
	}
	return nil, s.peer.Init(ctx, in.ID, in.State)
}

// done unmarshals the rpc request parameters and invokes the peer.Done
// procedure. The results are retuned and written to the rpc response.
func (s *Server) done(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	in := new(updateReq)
	if err := json.Unmarshal([]byte(*req.Params), in); err != nil {
		return nil, err
	}
	return nil, s.peer.Done(ctx, in.ID, in.State)
}

// extend unmarshals the rpc request parameters and invokes the peer.Extend
// procedure. The results are retuned and written to the rpc response.
func (s *Server) extend(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	var id string
	err := json.Unmarshal([]byte(*req.Params), &id)
	if err != nil {
		return nil, err
	}
	return nil, s.peer.Extend(ctx, id)
}
