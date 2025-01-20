package tradovate

import (
	"context"
	"fmt"

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

func (s *Socket) readFrame(ctx context.Context) (frame, error) {
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
