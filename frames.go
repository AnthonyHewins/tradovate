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
		f, err := s.readFrame(ctx)
		if err != nil {
			status := websocket.CloseStatus(err)
			if status == -1 {
				status = websocket.StatusInternalError
			}

			s.ws.Close(status, err.Error())
			return
		}

		switch f.frameType() {
		case frameTypeClose:
			s.ws.Close(0, "close received")
		case frameTypeData:
			msgs := f.(dataframe).msgs
			for i := range msgs {
				s.fm.pub(&msgs[i])
			}
		case frameTypeHeartbeat:
			s.ping(ctx)
		case frameTypeOpen:
			t, err := s.rest.Token(ctx)
			if err != nil {
				return
			}
			s.do(ctx, accessTokenURL, nil, t.AccessToken, nil)
		default:

		}
	}
}

func (s *WS) ping(ctx context.Context) error {
	return s.ws.Write(ctx, websocket.MessageText, []byte("[]"))
}
