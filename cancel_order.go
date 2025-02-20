package tradovate

import "context"

const cancelOrderURL = "order/cancelorder"

func (s *WS) CancelOrder(ctx context.Context, orderID uint) (commandID uint, err error) {
	type cancelResp struct {
		Fail OrderErrReason `json:"failureReason"`
		Text string         `json:"failureText"`
		Cmd  uint           `json:"commandId"`
	}

	var x cancelResp
	if err = s.do(ctx, cancelOrderURL, nil, map[string]uint{"orderId": orderID}, &x); err != nil {
		return commandID, err
	}

	if x.Fail != 0 {
		return 0, &OrderErr{Reason: x.Fail, Text: x.Text}
	}

	return x.Cmd, nil
}
