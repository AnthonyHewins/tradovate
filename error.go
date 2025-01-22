package tradovate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Err struct {
	Status int
	Body   string
}

func (e *Err) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Status, e.Body)
}

func newErrFromResp(r *http.Response) error {
	var e Err
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return fmt.Errorf(
			"while trying to read error response (HTTP %d) got err %w",
			r.StatusCode,
			err,
		)
	}

	return &e
}

func newErrFromSocket(r *rawMsg) error {
	var errmsg string
	if err := json.Unmarshal(r.Data, &errmsg); err != nil {
		return fmt.Errorf(
			"while trying to parse socket err msg with status %d, failed with %w. Raw: %s",
			r.Status,
			err,
			r.Data,
		)
	}

	return &Err{Status: r.Status, Body: errmsg}
}
