package tradovate

import (
	"context"
	"time"

	"github.com/goccy/go-json"
)

const (
	subscribeHistogram   = "md/subscribeHistogram"
	unsubscribeHistogram = "md/unsubscribeHistogram"
)

type Histogram struct {
	ContractID int                `json:"contractId"`
	Timestamp  time.Time          `json:"timestamp"`
	TradeDate  time.Time          `json:"tradeDate"` // set to 00:00:00-0500 (NYSE timezone)
	Base       float64            `json:"base"`
	Items      map[string]float64 `json:"items"`
	Refresh    bool               `json:"refresh"`
}

func (h *Histogram) UnmarshalJSON(b []byte) error {
	type histogram struct {
		ContractID int                `json:"contractId"`
		Timestamp  time.Time          `json:"timestamp"`
		TradeDate  tradeDate          `json:"tradeDate"`
		Base       float64            `json:"base"`
		Items      map[string]float64 `json:"items"`
		Refresh    bool               `json:"refresh"`
	}

	var x histogram
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	*h = Histogram{
		ContractID: x.ContractID,
		Timestamp:  x.Timestamp,
		TradeDate:  x.TradeDate.time(),
		Base:       x.Base,
		Items:      x.Items,
		Refresh:    x.Refresh,
	}
	return nil
}

func (s *WS) SubscribeHistogramID(ctx context.Context, id int) error {
	return s.subscribeHistogram(ctx, id)
}

func (s *WS) SubscribeHistogramSymbol(ctx context.Context, symbol string) error {
	return s.subscribeHistogram(ctx, symbol)
}

func (s *WS) subscribeHistogram(ctx context.Context, x any) error {
	return s.do(ctx, subscribeHistogram, nil, map[string]any{"symbol": x}, nil)
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
