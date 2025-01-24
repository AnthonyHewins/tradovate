package tradovate

import (
	"context"
	"fmt"

	"github.com/coder/websocket"
	"github.com/goccy/go-json"
)

type frame interface{ frameType() frameType }

type frameType byte

const (
	frameTypeUndefined frameType = iota
	frameTypeOpen
	frameTypeClose
	frameTypeHeartbeat
	frameTypeData
)

func (f frameType) frameType() frameType { return f }

type dataframe struct{ msgs []rawMsg }

func (d dataframe) frameType() frameType { return frameTypeData }

func (s *WS) readFrame(ctx context.Context) (frame, error) {
	_, binary, err := s.ws.Read(ctx)
	if err != nil {
		return nil, err
	}

	switch b := string(binary); b[0] {
	case 'o':
		return frameTypeOpen, nil
	case 'h':
		return frameTypeHeartbeat, nil
	case 'c':
		return frameTypeClose, nil
	case 'a':
		// no-op, parse message
	default:
		return nil, fmt.Errorf("unknown frame type received: %s. Raw: %s", string(b[0]), b)
	}

	var slice []rawMsg
	if err = json.UnmarshalContext(ctx, binary[1:], &slice); err != nil {
		return nil, err
	}

	return dataframe{slice}, nil
}

type rawMsg struct {
	Event  string          `json:"e"`
	ID     int             `json:"i"`
	Status int             `json:"s"`
	Data   json.RawMessage `json:"d"`
}

func (s *WS) keepalive(ctx context.Context) {
	for {
		closeDetails, err := s.keepaliveEvent(ctx)
		if err != nil {
			s.errHandler(err)
		}
		if closeDetails.status != 0 {
			s.connCancel()
			s.ws.Close(closeDetails.status, closeDetails.msg)
			return
		}
	}
}

type closeWith struct {
	status websocket.StatusCode
	msg    string
}

func (s *WS) keepaliveEvent(ctx context.Context) (closeWith, error) {
	f, err := s.readFrame(ctx)
	if err != nil {
		status := websocket.CloseStatus(err)
		if status == -1 || status == 0 {
			status = websocket.StatusInternalError
		}

		return closeWith{status, err.Error()}, err
	}

	switch f.frameType() {
	case frameTypeClose:
		return closeWith{websocket.StatusNormalClosure, "server closed conn"}, nil
	case frameTypeData:
		msgs := f.(dataframe).msgs
		for i := range msgs {
			s.fm.pub(&msgs[i])
		}
	case frameTypeHeartbeat:
		if err = s.ping(ctx); err != nil {
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

func (s *WS) ping(ctx context.Context) (err error) {
	var i uint8 = 0
	for ; i < s.pingRetries; i++ {
		if err = s.ws.Write(ctx, websocket.MessageText, []byte("[]")); err == nil {
			return nil
		}
	}

	return err
}
