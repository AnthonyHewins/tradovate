package tradovate

import (
	"context"
	"errors"
	"fmt"

	"github.com/coder/websocket"
)

var (
	ErrForceShutdown = errors.New("shutdown frame received")
)

func (s *WS) keepalive(ctx context.Context) {
	for {
		f, err := s.readFrame(ctx)
		if err != nil {
			status := websocket.CloseStatus(err)
			if status == -1 || status == 0 {
				status = websocket.StatusInternalError
			}

			s.errHandler(err)
			s.ws.Close(status, err.Error())
			return
		}

		go func(ctx context.Context, f frame) {
			closeDetails, err := s.keepaliveEvent(ctx, f)
			if err != nil {
				s.errHandler(err)
			}
			if closeDetails.status != 0 {
				s.connCancel()
				s.ws.Close(closeDetails.status, closeDetails.msg)
				return
			}
		}(ctx, f)
	}
}

type closeWith struct {
	status websocket.StatusCode
	msg    string
}

func (s *WS) keepaliveEvent(ctx context.Context, f frame) (closeWith, error) {
	switch f.frameType() {
	case frameTypeClose:
		return closeWith{websocket.StatusNormalClosure, "server closed conn"}, nil
	case frameTypeData:
		msgs := f.(dataframe).msgs
		return s.handleDataframe(msgs)
	case frameTypeHeartbeat:
		if err := s.ping(ctx); err != nil {
			return closeWith{
				3008, // timeout
				fmt.Sprintf("failed to ping after %d attempts", s.pingRetries),
			}, err
		}
	case frameTypeOpen:
		t, err := s.rest.Token(ctx)
		if err != nil {
			return closeWith{3000, "could not fetch auth token"}, err
		}

		if err = s.do(ctx, accessTokenURL, nil, t.AccessToken, nil); err != nil {
			return closeWith{3000, "unable to send auth token over socket"}, err
		}
	default:
		return closeWith{
			websocket.StatusInternalError,
			"invalid frame code received",
		}, fmt.Errorf("invalid frame type: %v", f.frameType())
	}

	return closeWith{}, nil
}

func (s *WS) handleDataframe(msgs []rawMsg) (closeWith, error) {
	for i := range msgs {
		switch v := &msgs[i]; v.Event {
		case frameEventUnspecified: // server response to request
			s.fm.pub(v)
		case frameEventClock:
			// unimplemented rn
		case frameEventProps: // server event update
			return eventHandler(v.entityMsg, s.entityHandler)
		case frameEventChart:
			return eventHandler(v.chart, s.chartHandler)
		case frameEventMd:
			return eventHandler(v.marketData, s.marketDataHandler)
		case frameEventShutdown:
			x, err := v.shutdownMsg()
			if err != nil {
				return closeWith{status: websocket.StatusInternalError, msg: err.Error()}, err
			}
			return closeWith{status: websocket.StatusNormalClosure, msg: x.Reason}, x
		default:
			s := fmt.Sprintf("unknown event type received: %s", v.Event)
			return closeWith{status: websocket.StatusInternalError, msg: s}, errors.New(s)
		}
	}

	return closeWith{}, nil
}

func eventHandler[X any](fn func() (X, error), handler func(X)) (closeWith, error) {
	e, err := fn()
	if err != nil {
		return closeWith{status: websocket.StatusInternalError, msg: err.Error()}, err
	}
	handler(e)
	return closeWith{}, nil
}

func (s *WS) readFrame(ctx context.Context) (frame, error) {
	_, binary, err := s.ws.Read(ctx)
	if err != nil {
		return nil, err
	}

	return newFrame(binary)
}

func (s *WS) ping(ctx context.Context) (err error) {
	var i uint8 = 0
	for ; i < s.pingRetries; i++ {
		if err = s.ws.Write(ctx, websocket.MessageText, []byte("[]")); err == nil {
			return nil
		}
	}

	return err
}
