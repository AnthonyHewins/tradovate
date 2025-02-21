package tests

import (
	"testing"
	"time"

	"github.com/AnthonyHewins/tradovate"
)

func TestMarketdata(t *testing.T) {
	resp, err := c.md.GetChartSymbol(c.ctx, "ES", &tradovate.ChartReq{
		UnderlyingType:   tradovate.ChartTypeMinuteBar,
		ElementSize:      15,
		ElementSizeUnit:  tradovate.SizeUnitUnderlyingUnits,
		WithHistogram:    false,
		AsFarAsTimestamp: time.Date(2025, 02, 20, 21, 0, 0, 0, time.UTC),
	})

	if err != nil {
		t.Errorf("failed fetching chart data: %s", err)
		return
	}
	defer c.md.CancelChart(c.ctx, resp.HistoricalID)

	select {
	case <-c.ctx.Done():
		t.Errorf("ctx canceled before data received")
	case v := <-c.chartChannel:
		if v == nil {
			t.Error("should not have gotten nil in chart test from channel")
			return
		}
	}
}
