package tradovate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	Prod = "https://live.tradovateapi.com/v1"
)

//go:generate interfacer -for github.com/AnthonyHewins/tradovate.Client -as tradovate.RESTInterface -o interface.go
type Client struct {
	tokenManager
	baseURL string
	h       *http.Client
}

// Credentials for getting a token
type Creds struct {
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	AppID      string    `json:"appId"`
	AppVersion string    `json:"appVersion"`
	CID        string    `json:"cid"`
	DeviceID   uuid.UUID `json:"deviceId"`
	Sec        uuid.UUID `json:"sec"`
}

func NewClient(baseURL string, h *http.Client, o *Creds) *Client {
	return &Client{
		tokenManager: tokenManager{creds: o},
		h:            h,
	}
}

func (r *Client) do(ctx context.Context, method, path string, reqBody, target any) error {
	buf, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		r.baseURL+path,
		bytes.NewReader(buf),
	)

	if err != nil {
		return err
	}

	t, err := r.Token(ctx)
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("authorization", "Bearer "+t.AccessToken)

	resp, err := r.h.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return newErrFromResp(resp)
	}

	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed decoding resp: %w", err)
	}

	return nil
}
