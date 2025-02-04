package tradovate

import (
	"encoding/json"
)

//go:generate enumer -type EventType -json -trimprefix EventType
type EventType byte

const (
	EventTypeUnspecified EventType = iota
	EventTypeCreated
	EventTypeUpdated
	EventTypeDeleted
)

//go:generate enumer -type EntityType -trimprefix EntityMsg -json
type EntityType byte

// TODO this list is way too long but I have no list of entity types
// from their docs. So this is the best I can do to make it simple for
// the end user and not store giant strings but just a byte
const (
	EntityTypeUnspecified EntityType = iota
	EntityTypeAccount
	EntityTypeAccountRiskStatus
	EntityTypeAdminAlert
	EntityTypeAdminAlertSignal
	EntityTypeCashBalance
	EntityTypeCashBalanceLog
	EntityTypeChat
	EntityTypeChatMessage
	EntityTypeClearingHouse
	EntityTypeCommand
	EntityTypeCommandReport
	EntityTypeContactInfo
	EntityTypeContract
	EntityTypeContractGroup
	EntityTypeContractMargin
	EntityTypeContractMaturity
	EntityTypeCurrency
	EntityTypeCurrencyRate
	EntityTypeEntitlement
	EntityTypeExchange
	EntityTypeExecutionReport
	EntityTypeFill
	EntityTypeFillFee
	EntityTypeFillPair
	EntityTypeMarginSnapshot
	EntityTypeMarketDataSubscription
	EntityTypeMarketDataSubscriptionExchangeScope
	EntityTypeMarketDataSubscriptionPlan
	EntityTypeOrder
	EntityTypeOrderStrategy
	EntityTypeOrderStrategyLink
	EntityTypeOrderStrategyType
	EntityTypeOrderVersion
	EntityTypeOrganization
	EntityTypePermissionedAccountAutoLiq
	EntityTypePosition
	EntityTypeProduct
	EntityTypeProductMargin
	EntityTypeProductSession
	EntityTypeProperty
	EntityTypeSecondMarketDataSubscription
	EntityTypeSpreadDefinition
	EntityTypeTradingPermission
	EntityTypeTradovateSubscription
	EntityTypeTradovateSubscriptionPlan
	EntityTypeUser
	EntityTypeUserAccountAutoLiq
	EntityTypeUserAccountPositionLimit
	EntityTypeUserAccountRiskParameter
	EntityTypeUserPlugin
	EntityTypeUserProperty
	EntityTypeUserSession
	EntityTypeUserSessionStats
)

type EntityMsg struct {
	Event EventType       `json:"eventType"`
	Type  EntityType      `json:"entityType"`
	Data  json.RawMessage `json:"entity"`
}
