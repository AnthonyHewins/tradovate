package tradovate

import (
	"context"
	"time"
)

const (
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
	OrderStatusUnknown OrderStatus = iota
	OrderStatusCanceled
	OrderStatusCompleted
	OrderStatusExpired
	OrderStatusFilled
	OrderStatusPendingCancel
	OrderStatusPendingNew
	OrderStatusPendingReplace
	OrderStatusRejected
	OrderStatusSuspended
	OrderStatusWorking
)

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

func (e *EntityMsg) Order() (*Order, error) { return decode[*Order](e) }
func (e *EntityMsg) MustOrder() *Order {
	o, err := e.Order()
	if err != nil {
		panic(err)
	}

	return o
}

func (s *WS) ListOrders(ctx context.Context) ([]*Order, error) {
	var x []*Order
	if err := s.do(ctx, listOrdersPath, nil, nil, &x); err != nil {
		return nil, err
	}

	return x, nil
}
