package tradovate

// Represents a single packet of market data.
// Only one of these slices should ever be set at a given
// time, but there's no way to know ahead of time which it is.
// So this API exposes all at once
type MarketData struct {
	Quotes     []*Quote     `json:"quotes"`
	DOMs       []*DOM       `json:"doms"`
	Histograms []*Histogram `json:"histograms"`
}
