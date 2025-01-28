package tradovate

import (
	"context"
	"time"
)

const (
	subscribeDOMs   = "md/subscribeDOM"
	unsubscribeDOMs = "md/unsubscribeDOM"
)

type DOM struct {
	ContractID int        `json:"contractId"`
	Timestamp  time.Time  `json:"timestamp"`
	Bids       []PriceQty `json:"bids"`
	Offers     []PriceQty `json:"offers"`
}

func (s *WS) SubscribeDOMSymbol(ctx context.Context, symbol string) error {
	return s.subscribeDOM(ctx, symbol)
}

func (s *WS) SubscribeDOMID(ctx context.Context, id int) error {
	return s.subscribeDOM(ctx, id)
}

func (s *WS) subscribeDOM(ctx context.Context, x any) error {
	return s.do(ctx, subscribeDOMs, nil, map[string]any{"symbol": x}, nil)
}

func (s *WS) UnsubscribeDOMSymbol(ctx context.Context, symbol string) error {
	return s.unsubscribeDOM(ctx, symbol)
}

func (s *WS) UnsubscribeDOMID(ctx context.Context, id int) error {
	return s.unsubscribeDOM(ctx, id)
}

func (s *WS) unsubscribeDOM(ctx context.Context, x any) error {
	return s.do(ctx, unsubscribeDOMs, nil, map[string]any{"symbol": x}, nil)
}
