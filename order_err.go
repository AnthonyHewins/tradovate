package tradovate

import (
	"strings"
)

//go:generate enumer -type OrderErrReason -trimprefix OrderErrReason -json
type OrderErrReason byte

const (
	OrderErrReasonAccountUnspecified OrderErrReason = iota
	OrderErrReasonAccountClosed
	OrderErrReasonAdvancedTrailingStopUnsupported
	OrderErrReasonAnotherCommandPending
	OrderErrReasonBackMonthProhibited
	OrderErrReasonExecutionProviderNotConfigured
	OrderErrReasonExecutionProviderUnavailable
	OrderErrReasonInvalidContract
	OrderErrReasonInvalidPrice
	OrderErrReasonLiquidationOnly
	OrderErrReasonLiquidationOnlyBeforeExpiration
	OrderErrReasonMaxOrderQtyIsNotSpecified
	OrderErrReasonMaxOrderQtyLimitReached
	OrderErrReasonMaxPosLimitMisconfigured
	OrderErrReasonMaxPosLimitReached
	OrderErrReasonMaxTotalPosLimitReached
	OrderErrReasonMultipleAccountPlanRequired
	OrderErrReasonNoQuote
	OrderErrReasonNotEnoughLiquidity
	OrderErrReasonOtherExecutionRelated
	OrderErrReasonParentRejected
	OrderErrReasonRiskCheckTimeout
	OrderErrReasonSessionClosed
	OrderErrReasonSuccess
	OrderErrReasonTooLate
	OrderErrReasonTradingLocked
	OrderErrReasonTrailingStopNonOrderQtyModify
	OrderErrReasonUnauthorized
	OrderErrReasonUnknownReason
	OrderErrReasonUnsupported
)

type OrderErr struct {
	Reason OrderErrReason
	Text   string
}

func (o *OrderErr) Error() string {
	var sb strings.Builder
	sb.WriteString(o.Reason.String())

	if o.Text != "" {
		sb.WriteString(": " + o.Text)
	}

	return sb.String()
}
