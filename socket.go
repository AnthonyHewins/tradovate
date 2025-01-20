package tradovate

import (
	"context"
	"sync/atomic"

	"github.com/coder/websocket"
)

const (
	WSSSandboxURL = "wss://demo.tradovateapi.com/v1/websocket"
	WSSReplayURL  = "wss://replay.tradovateapi.com/v1/websocket"
)

//go:generate interfacer -for github.com/AnthonyHewins/tradovate.Socket -as tradovate.SocketInterface -o socket_interface.go
type Socket struct {
	idGen atomic.Int64
	ws    *websocket.Conn
	fanout
}

// Create a new socket that will authenticate, use default websocket dialing
// options, then return the connection
func NewDefaultSocket(ctx context.Context) (*Socket, error) {
	return NewSocket(ctx, WSSSandboxURL, nil)
}

type SocketOpts struct {
	DialOpts *websocket.DialOptions
}

// Create a new socket that will authenticate with extra dialing options
func NewSocket(ctx context.Context, uri string, opts *SocketOpts) (*Socket, error) {
	conn, _, err := websocket.Dial(ctx, uri, opts.DialOpts)
	if err != nil {
		return nil, err
	}

	s := &Socket{ws: conn}
	go s.keepalive(ctx)
	return s, nil
}

// Ping sends a ping message to tradovate in the format they expect (JSON literal "[]")
func (s *Socket) Ping(ctx context.Context) error {
	return s.ws.Write(ctx, websocket.MessageText, []byte("[]"))
}

func (s *Socket) Close(code websocket.StatusCode, reason string) error {
	return s.ws.Close(code, reason)
}

func (s *Socket) keepalive(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():

		}
	}
}
