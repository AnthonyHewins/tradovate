package tradovate

import (
	"context"
	"time"
)

const (
	placeOsoOrderPath = "order/placeoso"
)

type OsoReq struct {
	AccountSpec    string      `json:"accountSpec"`
	AccountID      uint        `json:"accountId"`
	ClOrdID        string      `json:"clOrdId"`
	Action         Action      `json:"action"`
	Symbol         string      `json:"symbol"`
	OrderQty       uint        `json:"orderQty"`
	OrderType      OrderType   `json:"orderType"`
	Price          float64     `json:"price"`
	StopPrice      float64     `json:"stopPrice"`
	MaxShow        uint32      `json:"maxShow"`
	PegDifference  float64     `json:"pegDifference"`
	TimeInForce    Tif         `json:"timeInForce"`
	ExpireTime     time.Time   `json:"expireTime"`
	Text           string      `json:"text"`
	ActivationTime time.Time   `json:"activationTime"`
	CustomTag50    string      `json:"customTag50"`
	IsAutomated    bool        `json:"isAutomated"`
	Bracket1       *OtherOrder `json:"bracket1"`
	Bracket2       *OtherOrder `json:"bracket2"`
}

type OsoResp struct {
	OrderID, Oso1ID, Oso2ID uint
}

func (s *WS) OSO(ctx context.Context, o *OsoReq) (*OsoResp, error) {
	type osoResp struct {
		FailReason OrderErrReason `json:"failureReason"`
		FailText   string         `json:"failureText"`
		OrderID    uint           `json:"orderId"`
		OsoID1     uint           `json:"osoId1"`
		OsoID2     uint           `json:"osoId2"`
	}

	var x osoResp
	if err := s.do(ctx, placeOsoOrderPath, nil, o, &x); err != nil {
		return nil, err
	}

	if x.FailReason == OrderErrReasonSuccess {
		return nil, &OrderErr{Reason: x.FailReason, Text: x.FailText}
	}

	return &OsoResp{OrderID: x.OrderID, Oso1ID: x.OsoID1, Oso2ID: x.OsoID2}, nil
}
