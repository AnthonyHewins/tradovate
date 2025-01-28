package tradovate

import (
	"context"
	"time"
)

const placeOrderPath = "order/placeorder"

type OrderReq struct {
	AccountSpec    string    `json:"accountSpec"` // <= 64 chars: account username
	AccountID      int       `json:"accountId"`
	ClOrdId        string    `json:"clOrdId"` // string <= 64 characters
	Action         Action    `json:"action"`
	Symbol         string    `json:"symbol"` // string <= 64 characters
	OrderQty       uint32    `json:"orderQty"`
	OrderType      OrderType `json:"orderType"`
	Price          float64   `json:"price"`
	StopPrice      float64   `json:"stopPrice"`
	MaxShow        uint32    `json:"maxShow"`
	PegDifference  float64   `json:"pegDifference"`
	TimeInForce    Tif       `json:"timeInForce"`
	ExpireTime     time.Time `json:"expireTime"`
	Text           string    `json:"text"`
	ActivationTime time.Time `json:"activationTime"`
	CustomTag50    string    `json:"customTag50"`
	IsAutomated    bool      `json:"isAutomated"`
}

func (s *WS) PlaceOrder(ctx context.Context, r *OrderReq) (orderID int, err error) {
	type orderResp struct {
		Err  OrderErrReason `json:"failureReason"`
		Text string         `json:"failureText"`
		ID   int            `json:"orderId"`
	}

	var o orderResp
	if err := s.do(ctx, placeOrderPath, nil, r, &o); err != nil {
		return 0, err
	}

	if o.Err == OrderErrReasonAccountUnspecified {
		return o.ID, nil
	}

	return 0, &OrderErr{Reason: o.Err, Text: o.Text}
}
