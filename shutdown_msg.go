package tradovate

import (
	"encoding/json"
	"strings"
)

//go:generate enumer -type ShutdownCode -trimprefix ShutdownCode -json
type ShutdownCode byte

const (
	ShutdownCodeUnspecified ShutdownCode = iota
	ShutdownCodeMaintenance
	ShutdownCodeConnectionQuotaReached
	ShutdownCodeIPQuotaReached
)

type ShutdownMsg struct {
	Code   ShutdownCode `json:"reasonCode"`
	Reason string       `json:"reason"`
}

func (s *ShutdownMsg) Error() string {
	var sb strings.Builder
	sb.WriteString("shutdown received with code " + s.Code.String())

	if s.Reason != "" {
		sb.WriteString(". Reason: " + s.Reason)
	}

	return sb.String()
}

func decode[X any](e *EntityMsg) (X, error) {
	var x X
	return x, json.Unmarshal(e.Data, x)
}
