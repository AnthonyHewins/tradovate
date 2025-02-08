package tradovate

import (
	"context"
	"time"

	"github.com/goccy/go-json"
)

const (
	positionListURL = "position/list"
)

var (
	nyseTimezone *time.Location
)

func init() {
	var err error
	nyseTimezone, err = time.LoadLocation("America/New_York")
	if err != nil {
		nyseTimezone = time.FixedZone("America/New_York", int(-60*60*5))
	}
}

type tradeDate struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

func (t tradeDate) time() time.Time {
	return time.Date(t.Year, time.Month(t.Month), t.Day, 0, 0, 0, 0, nyseTimezone)
}

type Position struct {
	ID          int       `json:"id"`
	AccountID   int       `json:"accountId"`
	ContractID  int       `json:"contractId"`
	Timestamp   time.Time `json:"timestamp"`
	TradeDate   time.Time `json:"tradeDate"` // set to 00:00:00-0500 (nyse timezone)
	NetPos      int       `json:"netPos"`
	NetPrice    float64   `json:"netPrice"`
	Bought      int       `json:"bought"`
	BoughtValue float64   `json:"boughtValue"`
	Sold        int       `json:"sold"`
	SoldValue   float64   `json:"soldValue"`
	PrevPos     int       `json:"prevPos"`
	PrevPrice   float64   `json:"prevPrice"`
}

func (p *Position) UnmarshalJSON(b []byte) error {
	type position struct {
		ID          int       `json:"id"`
		AccountID   int       `json:"accountId"`
		ContractID  int       `json:"contractId"`
		Timestamp   time.Time `json:"timestamp"`
		TradeDate   tradeDate `json:"tradeDate"` // set to 00:00:00-0500 (nyse timezone)
		NetPos      int       `json:"netPos"`
		NetPrice    float64   `json:"netPrice"`
		Bought      int       `json:"bought"`
		BoughtValue float64   `json:"boughtValue"`
		Sold        int       `json:"sold"`
		SoldValue   float64   `json:"soldValue"`
		PrevPos     int       `json:"prevPos"`
		PrevPrice   float64   `json:"prevPrice"`
	}

	var x position
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}

	*p = Position{
		ID:          x.ID,
		AccountID:   x.AccountID,
		ContractID:  x.ContractID,
		Timestamp:   x.Timestamp,
		TradeDate:   x.TradeDate.time(),
		NetPos:      x.NetPos,
		NetPrice:    x.NetPrice,
		Bought:      x.Bought,
		BoughtValue: x.BoughtValue,
		Sold:        x.Sold,
		SoldValue:   x.SoldValue,
		PrevPos:     x.PrevPos,
		PrevPrice:   x.PrevPrice,
	}
	return nil
}

func (e *EntityMsg) Position() (*Position, error) { return decode[*Position](e) }
func (e *EntityMsg) MustPosition() *Position {
	d, err := e.Position()
	if err != nil {
		panic(err)
	}

	return d
}

func (s *WS) ListPositions(ctx context.Context) ([]*Position, error) {
	var positions []*Position
	err := s.do(ctx, positionListURL, nil, nil, &positions)
	if err != nil {
		return nil, err
	}

	return positions, nil
}
