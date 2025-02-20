package tests

import (
	"testing"
	"time"

	"github.com/AnthonyHewins/tradovate"
)

func TestOrderWorkflow(t *testing.T) {
	x, err := c.ws.PlaceOrder(c.ctx, &tradovate.OrderReq{
		AccountSpec:   c.spec,
		AccountID:     c.id,
		ClientOrderID: "asdjoisad",
		Action:        tradovate.ActionBuy,
		Symbol:        "NQU2",
		OrderQty:      1,
		OrderType:     tradovate.OrderTypeLimit,
		Price:         0.1,
		TimeInForce:   tradovate.TifDay,
		ExpireTime:    time.Now().Add(30 * time.Second),
		Text:          "integration test order",
		IsAutomated:   true,
	})

	if err != nil {
		t.Errorf("failed placing order: %s", err)
		return
	}

	o, err := c.ws.ListOrders(c.ctx)
	if err != nil {
		t.Errorf("failed listing orders: %s", err)
		return
	}

	found := false
	for _, v := range o {
		if found = v.ID == x; found {
			break
		}
	}

	if !found {
		t.Errorf("expected to see mock order created in order workflow, but didn't, got %+v", o)
		return
	}

	//	if _, err = c.ws.CancelOrder(c.ctx, x); err != nil {
	//		t.Errorf("failed canceling order during test: %v", err)
	//		return
	//	}
}
