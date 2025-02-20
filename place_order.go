package tradovate

import (
	"context"
	"time"
)

const placeOrderPath = "order/placeorder"

type OrderReq struct {
	AccountSpec    string    `json:"accountSpec,omitzero"` // <= 64 chars: account username
	AccountID      int       `json:"accountId,omitzero"`
	ClientOrderID  string    `json:"clOrdId,omitzero"` // string <= 64 characters
	Action         Action    `json:"action"`
	Symbol         string    `json:"symbol"` // string <= 64 characters
	OrderQty       uint32    `json:"orderQty"`
	OrderType      OrderType `json:"orderType"`
	Price          float64   `json:"price,omitzero"`
	StopPrice      float64   `json:"stopPrice,omitzero"`
	MaxShow        uint32    `json:"maxShow,omitzero"`
	PegDifference  float64   `json:"pegDifference,omitzero"`
	TimeInForce    Tif       `json:"timeInForce,omitzero"`
	ExpireTime     time.Time `json:"expireTime,omitzero"`
	Text           string    `json:"text,omitzero"`
	ActivationTime time.Time `json:"activationTime,omitzero"`
	CustomTag50    string    `json:"customTag50,omitzero"`
	IsAutomated    bool      `json:"isAutomated,omitzero"`
}

func (s *WS) PlaceOrder(ctx context.Context, r *OrderReq) (orderID uint, err error) {
	type orderResp struct {
		Err  OrderErrReason `json:"failureReason"`
		Text string         `json:"failureText"`
		ID   uint           `json:"orderId"`
	}

	var o orderResp
	if err := s.do(ctx, placeOrderPath, nil, r, &o); err != nil {
		return 0, err
	}

	if o.Err == OrderErrReasonSuccess {
		return o.ID, nil
	}

	return 0, &OrderErr{Reason: o.Err, Text: o.Text}
}
