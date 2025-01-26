package tradovate

import (
	"context"
	"encoding/json"
	"time"
)

const (
	getChart    = "md/getChart"
	cancelChart = "md/cancelChart"

	timestampFmt = "2006-01-02T15:04:05Z"
)

//go:generate enumer -type ChartType -json -trimprefix ChartType
type ChartType byte

const (
	ChartTypeUnspecified ChartType = iota
	ChartTypeTick
	ChartTypeDailyBar
	ChartTypeMinuteBar
	ChartTypeCustom
	ChartTypeDOM
)

//go:generate enumer -type SizeUnit -json -trimprefix SizeUnit
type SizeUnit byte

const (
	SizeUnitUnspecified SizeUnit = iota
	SizeUnitVolume
	SizeUnitRange
	SizeUnitUnderlyingUnits
	SizeUnitMomentumRange
	SizeUnitPointAndFigure
	SizeUnitOFARange
)

// Chart request to start receiving chart data
type ChartReq struct {
	UnderlyingType  ChartType
	ElementSize     uint32
	ElementSizeUnit SizeUnit
	WithHistogram   bool

	// One of these fields must be marked to be valid
	ClosestTimestamp time.Time
	ClosestTickID    uint32
	AsFarAsTimestamp time.Time
	AsMuchAsElements uint32
}

// This response is sent back to the user when a chart subscription has been enabled.
// Keep track of the realtime ID to cancel it appropriately
type ChartResp struct {
	HistoricalID int `json:"historicalId"`
	RealtimeID   int `json:"realtimeId"`
}

type Bar struct {
	Timestamp   time.Time `json:"timestamp"`
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	UpVolume    float64   `json:"upVolume"`
	DownVolume  float64   `json:"downVolume"`
	UpTicks     float64   `json:"upTicks"`
	DownTicks   float64   `json:"downTicks"`
	BidVolume   float64   `json:"bidVolume"`
	OfferVolume float64   `json:"offerVolume"`
}

type Tick struct {
	ID               int     `json:"id"`
	RelativeTime     int     `json:"t"` // relative tick timestamp in milliseconds
	RelativePrice    int     `json:"p"` // relative price - must be added to base price
	Volume           int     `json:"s"` // tick volume
	RelativeBidPrice float64 `json:"b"`
	RelativeAskPrice float64 `json:"a"`
	BidSize          float64 `json:"bs"`
	AskSize          float64 `json:"as"`
}

type Chart struct {
	ID int       // ID matching historicalId or realtimeId in ChartResp
	Td time.Time // trade date, set to 00:00:00Z

	// bar chart data
	Bars []Bar

	// Tick chart fields
	EndOfHistory  bool      // if this bool is set: the socket has finished loading historical data, now it'll be live
	Source        string    // if tick chart: source of data
	BasePrice     int       // contract tick sizes, ticks are calculated relative to this number
	BaseTimestamp time.Time // base timestamp, ticks calculated relative to this time
	TickSize      float64   // tick size that was requested
	Ticks         []Tick
}

func (c *Chart) UnmarshalJSON(b []byte) error {
	type chart struct {
		ID   int   `json:"id"`
		Td   int   `json:"td"` // timestamp as an int. very interesting choice here
		Bars []Bar `json:"bars"`
	}

	var cc chart
	if err := json.Unmarshal(b, &cc); err != nil {
		return err
	}

	// parse the time as a time.Time, it comes in as YYYYMMDD
	*c = Chart{
		ID: cc.ID,
		Td: time.Date(
			(cc.Td / 10000),
			time.Month((cc.Td/100)%100),
			cc.Td%100,
			0, 0, 0, 0, time.UTC,
		),
		Bars: cc.Bars,
	}
	return nil
}

func (s *WS) GetChartSymbol(ctx context.Context, symbol string, r *ChartReq) (ChartResp, error) {
	return s.getChart(ctx, symbol, r)
}

func (s *WS) GetChartID(ctx context.Context, id int, r *ChartReq) (ChartResp, error) {
	return s.getChart(ctx, id, r)
}

func (s *WS) getChart(ctx context.Context, x any, r *ChartReq) (ChartResp, error) {
	type chartDesc struct {
		UnderlyingType  ChartType `json:"underlyingType"`
		ElementSize     uint32    `json:"elementSize"`
		ElementSizeUnit SizeUnit  `json:"elementSizeUnit"`
		WithHistogram   bool      `json:"withHistogram"`
	}

	type timeRange struct {
		ClosestTimestamp string `json:"closestTimestamp"`
		ClosestTickID    uint32 `json:"closestTickId"`
		AsFarAsTimestamp string `json:"asFarAsTimestamp"`
		AsMuchAsElements uint32 `json:"asMuchAsElements"`
	}

	type chart struct {
		Symbol           any       `json:"symbol"`
		ChartDescription chartDesc `json:"chartDescription"`
		TimeRange        timeRange `json:"timeRange"`
	}

	c := chart{
		Symbol: x,
		ChartDescription: chartDesc{
			UnderlyingType:  r.UnderlyingType,
			ElementSize:     r.ElementSize,
			ElementSizeUnit: r.ElementSizeUnit,
			WithHistogram:   r.WithHistogram,
		},
		TimeRange: timeRange{
			ClosestTimestamp: r.ClosestTimestamp.Format(timestampFmt),
			ClosestTickID:    r.ClosestTickID,
			AsFarAsTimestamp: r.AsFarAsTimestamp.Format(timestampFmt),
			AsMuchAsElements: r.AsMuchAsElements,
		},
	}

	var resp ChartResp
	return resp, s.do(ctx, getChart, nil, &c, &resp)
}

// Cancel a chart subscription given the historicalId from ChartResp
func (s *WS) CancelChart(ctx context.Context, id int) error {
	return s.do(ctx, cancelChart, nil, map[string]any{"subscriptionId": id}, nil)
}
