package tests

import (
	"testing"
	"time"

	"github.com/AnthonyHewins/tradovate"
)

func TestMarketdata(t *testing.T) {
	resp, err := c.ws.GetChartSymbol(c.ctx, "ESM7", &tradovate.ChartReq{
		UnderlyingType:   tradovate.ChartTypeMinuteBar,
		ElementSize:      15,
		ElementSizeUnit:  tradovate.SizeUnitUnderlyingUnits,
		WithHistogram:    false,
		ClosestTimestamp: time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC),
		AsMuchAsElements: 50,
	})

	if err != nil {
		t.Errorf("failed fetching chart data: %s", err)
		return
	}
	defer c.ws.CancelChart(c.ctx, resp.HistoricalID)

	data := []*tradovate.Chart{}
	for v := range c.chartChannel {
		data = append(data, v)
	}

	t.Errorf("not enough data points in chart stream: %+v", err)
	return
}
