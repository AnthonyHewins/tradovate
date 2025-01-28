package tradovate

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/goccy/go-json"
)

const (
	WSSSandboxURL = "wss://demo.tradovateapi.com/v1/websocket"
	WSSReplayURL  = "wss://replay.tradovateapi.com/v1/websocket"
)

type WSOpt func(s *WS)

// Seed the client with a token if you have it persisted.
// Skips authentication at the beginning if not expired
func WithToken(t *Token) WSOpt {
	return func(s *WS) { s.rest.SetToken(t) }
}

// Attaches a timeout to each WS request. Defaults to 5s
// if you don't set
func WithTimeout(t time.Duration) WSOpt {
	return func(s *WS) { s.fm.deadline = t }
}

func WithErrHandler(fn func(error)) WSOpt {
	return func(s *WS) { s.errHandler = fn }
}

// How many times to retry pings until the websocket
// will assume a dead connection. Defaults to 5
func WithPingRetries(x uint8) WSOpt {
	return func(s *WS) {
		if x == 0 {
			x++
		}
		s.pingRetries = x
	}
}

// All entity events the server propagates that don't have to
// do with request-response will be called here
//
// Each call will run as a goroutine
func WithEntityHandler(x func(*EntityMsg)) WSOpt {
	return func(s *WS) { s.entityHandler = x }
}

// All chart messages will be received by this handler
//
// Each call will run as a goroutine
func WithChartHandler(x func(*Chart)) WSOpt {
	return func(s *WS) { s.chartHandler = x }
}

// Handle market data
//
// Each call will run as a goroutine
func WithMarketDataHandler(x func(*MarketData)) WSOpt {
	return func(s *WS) { s.marketDataHandler = x }
}

// Websocket client to the tradovate API
type WS struct {
	connCtx    context.Context
	connCancel context.CancelFunc

	pingRetries uint8
	rest        *REST
	ws          *websocket.Conn
	fm          fanoutMutex

	entityHandler     func(*EntityMsg)
	chartHandler      func(*Chart)
	marketDataHandler func(*MarketData)
	errHandler        func(error)
}

func NewSocket(ctx context.Context, uri string, dialOpts *websocket.DialOptions, rest *REST, opts ...WSOpt) (*WS, error) {
	if rest == nil {
		return nil, fmt.Errorf("missing rest client: need it for auth")
	}

	ws, _, err := websocket.Dial(ctx, uri, dialOpts)
	if err != nil {
		return nil, err
	}

	connCtx, cancel := context.WithCancel(ctx)

	s := &WS{
		pingRetries: 5,
		connCtx:     connCtx,
		connCancel:  cancel,
		rest:        rest,
		ws:          ws,
		fm: fanoutMutex{
			acc:      1,
			deadline: time.Second * 5,
		},
		entityHandler:     func(em *EntityMsg) {},
		chartHandler:      func(cr *Chart) {},
		marketDataHandler: func(md *MarketData) {},
		errHandler:        func(err error) {},
	}

	for _, v := range opts {
		v(s)
	}

	go s.keepalive(connCtx)
	return s, nil
}

func (s *WS) Close(ctx context.Context) error {
	return s.ws.Close(websocket.StatusNormalClosure, "client initiated close")
}

func (s *WS) do(ctx context.Context, path string, queryParams url.Values, body, target any) error {
	sb := strings.Builder{}

	sb.WriteString(path)
	sb.WriteRune('\n')

	mu := s.fm.request()
	sb.WriteString(fmt.Sprint(mu.id))
	sb.WriteRune('\n')

	if len(queryParams) > 0 {
		sb.WriteString(queryParams.Encode())
	}
	sb.WriteRune('\n')

	if body != nil {
		err := json.NewEncoder(&sb).EncodeContext(ctx, body)
		if err != nil {
			return err
		}
	}

	err := s.ws.Write(ctx, websocket.MessageText, []byte(sb.String()))
	if err != nil {
		return err
	}

	resp, err := mu.wait(ctx, s.connCtx)
	if err != nil {
		return err
	}

	if resp.Status >= 300 {
		return newRespErrFromSocket(resp)
	}

	if target != nil {
		return json.UnmarshalContext(ctx, resp.Data, target)
	}

	return nil
}
