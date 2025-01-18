package tradovate

import (
	"context"

	"github.com/coder/websocket"
)

const (
	SandboxURL = "wss://demo.tradovateapi.com/v1/websocket"
)

type Socket struct {
	ws *websocket.Conn
}

// Create a new socket that will authenticate, use default websocket dialing
// options, then return the connection
func NewDefaultSocket(ctx context.Context) (*Socket, error) {
	return NewSocket(ctx, SandboxURL, nil)
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

	return &Socket{ws: conn}, nil
}
