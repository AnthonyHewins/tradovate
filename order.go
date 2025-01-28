package tradovate

import (
	"context"
	"fmt"
	"time"
)

const (
	placeOrderPath = "order/placeorder"
	listOrdersPath = "order/list"
)

//go:generate enumer -type Action -trimprefix Action -json
type Action byte

const (
	ActionUnspecified Action = iota
	ActionBuy
	ActionSell
)

//go:generate enumer -type OrderType -trimprefix OrderType -json
type OrderType byte

const (
	OrderTypeUnspecified OrderType = iota
	OrderTypeLimit
	OrderTypeMIT
	OrderTypeMarket
	OrderTypeQTS
	OrderTypeStop
	OrderTypeStopLimit
	OrderTypeTrailingStop
	OrderTypeTrailingStopLimit
)

// Time in force
//
//go:generate enumer -type Tif -trimprefix TIF
type Tif byte

const (
	TifUnspecified Tif = iota
	TifDay
	TifFOK
	TifGTC
	TifGTD
	TifIOC
)

//go:generate enumer -type OrderStatus -trimprefix OrderStatus -json
type OrderStatus byte

const (
	OrderStatusUnspecified OrderStatus = iota
	OrderStatusCanceled
	OrderStatusCompleted
	OrderStatusExpired
	OrderStatusFilled
	OrderStatusPendingCancel
	OrderStatusPendingNew
	OrderStatusPendingReplace
	OrderStatusRejected
	OrderStatusSuspended
	OrderStatusUnknown
	OrderStatusWorking
)

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

type Order struct {
	ID                  int         `json:"id"`
	AccountID           uint        `json:"accountId"`
	ContractID          uint        `json:"contractId"`
	SpreadDefinitionID  uint        `json:"spreadDefinitionId"`
	Timestamp           time.Time   `json:"timestamp"`
	Action              Action      `json:"action"`
	Status              OrderStatus `json:"ordStatus"`
	ExecutionProviderID uint        `json:"executionProviderId"`
	OcoID               uint        `json:"ocoId"`
	ParentID            uint        `json:"parentId"`
	LinkedID            uint        `json:"linkedId"`
	Admin               bool        `json:"admin"`
}

func (s *WS) PlaceOrder(ctx context.Context, r *OrderReq) (orderID int, err error) {
	type orderResp struct {
		Err  string `json:"failureReason"`
		Text string `json:"failureText"`
		ID   int    `json:"orderId"`
	}

	var o orderResp
	if err := s.do(ctx, placeOrderPath, nil, r, &o); err != nil {
		return 0, err
	}

	if o.Err == "" {
		return o.ID, nil
	}

	return 0, fmt.Errorf("%s: %s", o.Err, o.Text)
}

func (s *WS) ListOrders(ctx context.Context) ([]*Order, error) {
	var x []*Order
	if err := s.do(ctx, listOrdersPath, nil, nil, &x); err != nil {
		return nil, err
	}

	return x, nil
}
