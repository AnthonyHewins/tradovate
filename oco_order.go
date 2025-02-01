package tradovate

import (
	"context"
	"time"
)

const (
	placeOcoOrder = "order/placeoco"
)

// The order in an OCO order, dubbed "other order" in tradovate's
// docs
type OtherOrder struct {
	Action        Action    `json:"action"`
	ClOrdID       string    `json:"clOrdId"`
	OrderType     OrderType `json:"orderType"`
	Price         float64   `json:"price"`
	StopPrice     float64   `json:"stopPrice"`
	MaxShow       uint32    `json:"maxShow"`
	PegDifference float64   `json:"pegDifference"`
	TimeInForce   Tif       `json:"timeInForce"`
	ExpireTime    time.Time `json:"expireTime"`
	Text          string    `json:"text"`
}

// One-cancels-other order
type OcoReq struct {
	AccountSpec    string     `json:"accountSpec"`
	AccountID      uint       `json:"accountId"`
	ClOrdID        string     `json:"clOrdId"`
	Action         Action     `json:"action"`
	Symbol         string     `json:"symbol"`
	OrderQty       uint       `json:"orderQty"`
	OrderType      OrderType  `json:"orderType"`
	Price          float64    `json:"price"`
	StopPrice      float64    `json:"stopPrice"`
	MaxShow        uint32     `json:"maxShow"`
	PegDifference  float64    `json:"pegDifference"`
	TimeInForce    Tif        `json:"timeInForce"`
	ExpireTime     time.Time  `json:"expireTime"`
	Text           string     `json:"text"`
	ActivationTime time.Time  `json:"activationTime"`
	CustomTag50    string     `json:"customTag50"`
	IsAutomated    bool       `json:"isAutomated"`
	Other          OtherOrder `json:"other"`
}

type OcoResp struct {
	OrderID, OcoID uint
}

func (s *WS) OCO(ctx context.Context, o *OcoReq) (*OcoResp, error) {
	type ocoResp struct {
		FailReason OrderErrReason `json:"failureReason"`
		FailText   string         `json:"failureText"`
		OrderID    uint           `json:"orderId"`
		OcoID      uint           `json:"ocoId"`
	}

	var x ocoResp
	if err := s.do(ctx, placeOcoOrder, nil, o, &x); err != nil {
		return nil, err
	}

	if x.FailReason == OrderErrReasonSuccess {
		return &OcoResp{OrderID: x.OrderID, OcoID: x.OcoID}, nil
	}

	return nil, &OrderErr{Reason: x.FailReason, Text: x.FailText}
}
