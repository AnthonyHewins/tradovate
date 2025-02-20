package tradovate

import (
	"context"
	"time"
)

const accountListURL = "account/list"

//go:generate enumer -type AccountType -json -trimprefix AccountType
type AccountType byte

const (
	AccountTypeUnspecified AccountType = iota
	AccountTypeCustomer
	AccountTypeGiveup
	AccountTypeHouse
	AccountTypeOmnibus
	AccountTypeWash
)

//go:generate enumer -type MarginAccount -trimprefix MarginAccount -json
type MarginAccount byte

const (
	MarginAccountUnspecified MarginAccount = iota
	MarginAccountHedger
	MarginAccountSpeculator
)

//go:generate enumer -type LegalStatus -trimprefix LegalStatus -json
type LegalStatus byte

const (
	LegalStatusUnspecified LegalStatus = iota
	LegalStatusCorporation
	LegalStatusGP
	LegalStatusIRA
	LegalStatusIndividual
	LegalStatusJoint
	LegalStatusLLC
	LegalStatusLLP
	LegalStatusLP
	LegalStatusPTR
	LegalStatusTrust
)

type Account struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	UserID            uint          `json:"userId"`
	AccountType       AccountType   `json:"accountType"`
	Active            bool          `json:"active"`
	ClearingHouseID   uint          `json:"clearingHouseId"`
	RiskCategoryID    uint          `json:"riskCategoryId"`
	AutoLiqProfileID  uint          `json:"autoLiqProfileId"`
	MarginAccountType MarginAccount `json:"marginAccountType"`
	LegalStatus       string        `json:"legalStatus"`
	Timestamp         time.Time     `json:"timestamp"`
	EvaluationSize    float64       `json:"evaluationSize"`
	Readonly          bool          `json:"readonly"`
	CcEmail           string        `json:"ccEmail"`
}

func (s *WS) ListAccounts(ctx context.Context) ([]*Account, error) {
	var x []*Account
	if err := s.do(ctx, accountListURL, nil, nil, &x); err != nil {
		return nil, err
	}
	return x, nil
}
