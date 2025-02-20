package tradovate

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"encoding/json"

	"github.com/coder/websocket"
)

const (
	WSSMarketDataURL = "wss://md.tradovateapi.com/v1/websocket"
	WSSSandboxURL    = "wss://demo.tradovateapi.com/v1/websocket"
	WSSLiveURL       = "wss://live.tradovateapi.com/v1/websocket"
	WSSReplayURL     = "wss://replay.tradovateapi.com/v1/websocket"
)

type WSOpt func(s *WS)

// Attaches a timeout to each WS request. Defaults to 5s
// if you don't set
func WithTimeout(t time.Duration) WSOpt {
	return func(s *WS) { s.fm.timeout = t }
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
	if uri == "" {
		return nil, fmt.Errorf("missing connection URI")
	}

	if rest == nil {
		return nil, fmt.Errorf("rest client nil: need rest client to authenticate")
	}

	connCtx, cancel := context.WithCancel(ctx)

	s := &WS{
		pingRetries: 5,
		connCtx:     connCtx,
		connCancel:  cancel,
		rest:        rest,
		fm: fanoutMutex{
			acc:     1,
			timeout: time.Second * 5,
		},
		entityHandler:     func(em *EntityMsg) {},
		chartHandler:      func(cr *Chart) {},
		marketDataHandler: func(md *MarketData) {},
		errHandler:        func(err error) {},
	}

	for _, v := range opts {
		v(s)
	}

	t, err := rest.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting initial token: %w", err)
	}

	s.ws, _, err = websocket.Dial(ctx, uri, dialOpts)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			s.ws.Close(websocket.StatusInternalError, "failed initial setup: "+err.Error())
		}
	}()

	f, err := s.readFrame(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed reading opening handshake packet with %s: %w", uri, err)
	}

	if f.frameType() != frameTypeOpen {
		return nil, fmt.Errorf("protocol broken: frame type should be open, but got %+v", f)
	}

	go s.keepalive(connCtx)
	return s, s.do(ctx, "authorize", nil, t.AccessToken, nil)
}

func (s *WS) Close() error {
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
		err := json.NewEncoder(&sb).Encode(body)
		if err != nil {
			return err
		}
	}

	payload := []byte(sb.String())
	fmt.Println(string(payload))
	if err := s.ws.Write(ctx, websocket.MessageText, payload); err != nil {
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
		return json.Unmarshal(resp.Data, target)
	}

	return nil
}
