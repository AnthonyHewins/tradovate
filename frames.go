package tradovate

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrEmptyFrame = errors.New("empty frame received")

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

func newFrame(binary []byte) (frame, error) {
	if len(binary) == 0 {
		return nil, ErrEmptyFrame
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
	if err := json.Unmarshal(binary[1:], &slice); err != nil {
		return nil, err
	}

	return dataframe{slice}, nil
}
