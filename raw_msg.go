package tradovate

import "encoding/json"

//go:generate enumer -type frameEvent -json -trimprefix frameEvent -transform lower
type frameEvent byte

const (
	frameEventUnspecified frameEvent = iota
	frameEventProps
	frameEventShutdown
	frameEventMd
	frameEventChart
	frameEventClock
)

type rawMsg struct {
	// The "e" field specifies an event kind:
	//
	//  - "props": this is a notification that some entity was created, updated or deleted. "d" field specifies details of the event with the next structure:
	//      -  "entityType" field
	//      -  "entity" field. JSON structure of object (or array of objects) specified in this field is identical to JSON of entity that accessible via corresponding REST API request like entityType/item. For example, if entityType=account, JSON can be found in the response specification of account/item call
	//      -  "eventType" field with options "Created", "Updated" or "Deleted"
	//  - "shutdown": a notification before graceful shutdown of connection. "d" field specifies details:
	//       "reasonCode" field with options "Maintenance", "ConnectionQuotaReached", "IPQuotaReached"
	//       "reason" field is optional and may contain a readable explanation
	//  - "md" and "chart": these notifications are used by market data feed services only, the description of "d" field is here
	//  - "clock": Market Replay clock synchronization message. See the Market Replay section below.
	Event  frameEvent      `json:"e"`
	ID     int             `json:"i"`
	Status int             `json:"s"`
	Data   json.RawMessage `json:"d"`
}

func (r *rawMsg) entityMsg() (*EntityMsg, error) {
	var e EntityMsg
	if err := json.Unmarshal(r.Data, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *rawMsg) chart() (*Chart, error) {
	var c Chart
	if err := json.Unmarshal(r.Data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *rawMsg) marketData() (*MarketData, error) {
	var md MarketData
	if err := json.Unmarshal(r.Data, &md); err != nil {
		return nil, err
	}

	return &md, nil
}

func (r *rawMsg) shutdownMsg() (*ShutdownMsg, error) {
	var s ShutdownMsg
	if err := json.Unmarshal(r.Data, &s); err != nil {
		return nil, err
	}

	return &s, nil
}
