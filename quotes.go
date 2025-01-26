package tradovate

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	subscribeQuotePath   = "md/subscribeQuote"
	unsubscribeQuotePath = "md/unsubscribeQuote"
)

type PriceQty struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

type Quote struct {
	ContractID int
	Timestamp  time.Time

	Bid   PriceQty
	Offer PriceQty
	Trade PriceQty

	TotalTradeVolume float64
	OpenInterest     float64

	LowPrice        float64
	OpeningPrice    float64
	HighPrice       float64
	SettlementPrice float64
}

// Implemented unmarshal for this data structure because it's
// a total mess of garbage fluff coming from the server. Flattened out
// it's much easier to reason with
func (q *Quote) UnmarshalJSON(b []byte) error {
	type size struct {
		Size float64 `json:"size"`
	}

	type price struct {
		Price float64 `json:"price"`
	}

	type bloat struct {
		Timestamp  time.Time `json:"timestamp"`
		ContractID int       `json:"contractId"`
		Entries    struct {
			Bid   PriceQty `json:"Bid"`
			Offer PriceQty `json:"Offer"`
			Trade PriceQty `json:"Trade"`

			TotalTradeVolume size `json:"TotalTradeVolume"`
			OpenInterest     size `json:"OpenInterest"`

			OpeningPrice    price `json:"OpeningPrice"`
			LowPrice        price `json:"LowPrice"`
			HighPrice       price `json:"HighPrice"`
			SettlementPrice price `json:"SettlementPrice"`
		} `json:"entries"`
	}

	var x bloat
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	e := x.Entries
	*q = Quote{
		Timestamp:        x.Timestamp,
		ContractID:       x.ContractID,
		Bid:              e.Bid,
		TotalTradeVolume: e.TotalTradeVolume.Size,
		Offer:            e.Offer,
		LowPrice:         e.LowPrice.Price,
		Trade:            e.Trade,
		OpenInterest:     e.OpenInterest.Size,
		OpeningPrice:     e.OpeningPrice.Price,
		HighPrice:        e.HighPrice.Price,
		SettlementPrice:  e.SettlementPrice.Price,
	}

	return nil
}

// Subscribe to a contract by ID. If you prefer doing it by contract ID, use
// SubscribeQuoteID
func (s *WS) SubscribeQuoteSymbol(ctx context.Context, symbol string) ([]*Quote, error) {
	if symbol == "" {
		return nil, fmt.Errorf("no symbol passed to subscribe")
	}

	return s.marketDataSubscribeQuote(ctx, symbol)
}

// Subscribe to a contract by ID. If you prefer doing it by symbol, use
// SubscribeQuoteSymbol
func (s *WS) SubscribeQuoteID(ctx context.Context, id int) ([]*Quote, error) {
	return s.marketDataSubscribeQuote(ctx, id)
}

// Subscribe to a contract by ID. If you prefer doing it by contract ID, use
// SubscribeQuoteID
func (s *WS) UnsubscribeQuoteSymbol(ctx context.Context, symbol string) error {
	if symbol == "" {
		return fmt.Errorf("no symbol passed to unsubscribe")
	}

	return s.unsubscribeQuote(ctx, symbol)
}

// Subscribe to a contract by ID. If you prefer doing it by symbol, use
// SubscribeQuoteSymbol
func (s *WS) UnsubscribeQuoteID(ctx context.Context, id int) error {
	return s.unsubscribeQuote(ctx, id)
}

func (s *WS) unsubscribeQuote(ctx context.Context, x any) error {
	return s.do(ctx, unsubscribeQuotePath, nil, map[string]any{"symbol": x}, nil)
}

func (s *WS) marketDataSubscribeQuote(ctx context.Context, x any) ([]*Quote, error) {
	var q []*Quote
	if err := s.do(ctx, subscribeQuotePath, nil, map[string]any{"symbol": x}, q); err != nil {
		return nil, err
	}
	return q, nil
}
