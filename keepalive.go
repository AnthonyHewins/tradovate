package tradovate

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/coder/websocket"
)

var (
	ErrForceShutdown = errors.New("shutdown frame received")
)

func (s *WS) closeErr(err error) {
	s.connCancel()
	go s.errHandler(err)

	status := websocket.CloseStatus(err)
	if status == -1 || status == 0 {
		status = websocket.StatusInternalError
	}

	s.ws.Close(status, err.Error())
}

func (s *WS) keepalive(ctx context.Context) {
	for {
		f, err := s.readFrame(ctx)
		if err != nil {
			switch {
			case errors.Is(err, context.Canceled):
				s.Close()
				return
			case err == net.ErrClosed:
				s.connCancel()
				s.errHandler(err)
				return
			}

			s.closeErr(err)
			return
		}

		switch f.frameType() {
		case frameTypeClose:
			err = websocket.CloseError{
				Code:   websocket.StatusAbnormalClosure,
				Reason: "unexpected close",
			}
		case frameTypeData:
			err = s.handleDataframe(f.(dataframe).msgs)
		case frameTypeHeartbeat:
			go func() { // dont slow down the read routine for ping writes
				if pingErr := s.ping(ctx); pingErr != nil {
					s.closeErr(pingErr)
				}
			}()
		case frameTypeOpen:
			var t *Token
			if t, err = s.rest.Token(ctx); err == nil {
				err = s.do(ctx, accessTokenURL, nil, t.AccessToken, nil)
			}
		default:
			err = websocket.CloseError{
				Code:   websocket.StatusInternalError,
				Reason: fmt.Sprintf("invalid frame type: %v", f.frameType()),
			}
		}

		if err != nil {
			s.closeErr(err)
		}
	}
}

func (s *WS) handleDataframe(msgs []rawMsg) error {
	for i := range msgs {
		var err error
		switch v := &msgs[i]; v.Event {
		case frameEventUnspecified: // server response to request
			s.fm.pub(v)
		case frameEventClock:
			// unimplemented rn
		case frameEventProps: // server event update
			err = eventHandler(v.entityMsg, s.entityHandler)
		case frameEventChart:
			err = eventHandler(v.chart, s.chartHandler)
		case frameEventMd:
			err = eventHandler(v.marketData, s.marketDataHandler)
		case frameEventShutdown:
			var x *ShutdownMsg
			if x, err = v.shutdownMsg(); err == nil {
				err = websocket.CloseError{Code: websocket.StatusNormalClosure, Reason: x.Reason}
			}
		default:
			return websocket.CloseError{
				Code:   websocket.StatusInternalError,
				Reason: fmt.Sprintf("unknown event type received: %s", v.Event),
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func eventHandler[X any](fn func() (X, error), handler func(X)) error {
	e, err := fn()
	if err == nil {
		go handler(e)
	}
	return err
}

func (s *WS) readFrame(ctx context.Context) (frame, error) {
	_, binary, err := s.ws.Read(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(binary))
	return newFrame(binary)
}

func (s *WS) ping(ctx context.Context) error {
	var i uint8 = 0
	for ; i < s.pingRetries; i++ {
		switch err := s.ws.Write(ctx, websocket.MessageText, []byte("[]")); {
		case err == nil || err == net.ErrClosed:
			return err
		case errors.Is(err, context.Canceled):
			return nil // main keepalive will error out
		}
	}

	return websocket.CloseError{
		Code:   3008, // timeout
		Reason: fmt.Sprintf("failed to ping after %d attempts", s.pingRetries),
	}
}
