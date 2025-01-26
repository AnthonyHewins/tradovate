package tradovate

import (
	"context"
	"time"
)

const (
	subscribeHistogram   = "md/subscribeHistogram"
	unsubscribeHistogram = "md/unsubscribeHistogram"
)

type Histogram struct {
	ContractID int                `json:"contractId"`
	Timestamp  time.Time          `json:"timestamp"`
	TradeDate  TradeDate          `json:"tradeDate"`
	Base       float64            `json:"base"`
	Items      map[string]float64 `json:"items"`
	Refresh    bool               `json:"refresh"`
}

func (s *WS) SubscribeHistogramID(ctx context.Context, id int) ([]*Histogram, error) {
	return s.subscribeHistogram(ctx, id)
}

func (s *WS) SubscribeHistogramSymbol(ctx context.Context, symbol string) ([]*Histogram, error) {
	return s.subscribeHistogram(ctx, symbol)
}

func (s *WS) subscribeHistogram(ctx context.Context, x any) ([]*Histogram, error) {
	type histogram struct {
		H []*Histogram `json:"histograms"`
	}

	var h histogram
	if err := s.do(ctx, subscribeHistogram, nil, map[string]any{"symbol": x}, &h); err != nil {
		return nil, err
	}

	return h.H, nil
}

func (s *WS) UnsubscribeHistogramID(ctx context.Context, id int) error {
	return s.unsubscribeHistogram(ctx, id)
}

func (s *WS) UnsubscribeHistogramSymbol(ctx context.Context, symbol string) error {
	return s.unsubscribeHistogram(ctx, symbol)
}

func (s *WS) unsubscribeHistogram(ctx context.Context, x any) error {
	return s.do(ctx, unsubscribeHistogram, nil, map[string]any{"symbol": x}, nil)
}
