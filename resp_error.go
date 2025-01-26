package tradovate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RespErr struct {
	Status int
	Body   string
}

func (e *RespErr) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Status, e.Body)
}

func (e *RespErr) Is(err error) bool {
	_, ok := err.(*RespErr)
	return ok
}

func newRespErrFromREST(r *http.Response) error {
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf(
			"while trying to read error response (HTTP %d) got err %w",
			r.StatusCode,
			err,
		)
	}
	r.Body.Close()

	return &RespErr{Status: r.StatusCode, Body: string(buf)}
}

func newRespErrFromSocket(r *rawMsg) error {
	var errmsg string
	if err := json.Unmarshal(r.Data, &errmsg); err != nil {
		return fmt.Errorf(
			"while trying to parse socket err msg with status %d, failed with %w. Raw: %s",
			r.Status,
			err,
			r.Data,
		)
	}

	return &RespErr{Status: r.Status, Body: errmsg}
}
