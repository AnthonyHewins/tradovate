package tradovate

import (
	"context"
	"time"
)

const (
	positionListURL = "position/list"
)

type TradeDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type Position struct {
	ID          int       `json:"id"`
	AccountID   int       `json:"accountId"`
	ContractID  int       `json:"contractId"`
	Timestamp   time.Time `json:"timestamp"`
	TradeDate   TradeDate `json:"tradeDate"`
	NetPos      int       `json:"netPos"`
	NetPrice    float64   `json:"netPrice"`
	Bought      int       `json:"bought"`
	BoughtValue float64   `json:"boughtValue"`
	Sold        int       `json:"sold"`
	SoldValue   float64   `json:"soldValue"`
	PrevPos     int       `json:"prevPos"`
	PrevPrice   float64   `json:"prevPrice"`
}

func (s *WS) ListPositions(ctx context.Context) ([]*Position, error) {
	var positions []*Position
	err := s.do(ctx, positionListURL, nil, nil, &positions)
	if err != nil {
		return nil, err
	}

	return positions, nil
}
