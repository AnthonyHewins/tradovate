package tradovate

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

const (
	accessTokenURL = "auth/accessTokenRequest"
)

// Token is the token response from fetching access tokens.
// Tokens have a lifespan of 90 minutes
type Token struct {
	AccessToken            string    `json:"accessToken"`
	ExpirationTime         time.Time `json:"expirationTime"`
	PasswordExpirationTime time.Time `json:"passwordExpirationTime"`
	UserStatus             string    `json:"userStatus"`
	UserID                 int       `json:"userId"`
	Name                   string    `json:"name"`
	HasLive                bool      `json:"hasLive"`
}

// Check if token is expired
func (t *Token) Expired() bool {
	return time.Now().After(t.ExpirationTime)
}

type tokenManager struct {
	mu                   sync.RWMutex
	forceRefreshDeadline time.Duration
	creds                *Creds
	token                *Token
}

// Sets the token in case you have it persisted somewhere else
// and want to seed it
func (t *tokenManager) SetToken(token *Token) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.token = token
}

// Fetches a token using the following steps:
//  1. If no token is set: return a new, fresh token (you can
//     avoid this by SetToken if you already have one)
//  2. If a token exists but is about to expire: fetch a fresh one
//  3. If a token exists but expires within the window you specified
//     when you created the client, it will refresh the token
//  4. Return the token otherwise, because it exists and isn't going to
//     expire soon
func (r *REST) Token(ctx context.Context) (*Token, error) {
	r.tokenManager.mu.RLock()
	x := r.tokenManager.token
	if x == nil || time.Now().Add(time.Millisecond*20).After(x.ExpirationTime) {
		r.tokenManager.mu.RUnlock()
		return r.newToken(ctx)
	}

	if time.Until(r.token.ExpirationTime) < r.tokenManager.forceRefreshDeadline {
		r.tokenManager.mu.RUnlock()
		return r.refreshToken(ctx)
	}

	defer r.mu.RUnlock()
	return x, nil
}

// this struct contains 1 extra field for error text.
// terrible API design IMO, but here we are. I hid this
// from the end user to simplify the return
type tokenResp struct {
	ErrorText              string    `json:"errorText"`
	AccessToken            string    `json:"accessToken"`
	ExpirationTime         time.Time `json:"expirationTime"`
	PasswordExpirationTime time.Time `json:"passwordExpirationTime"`
	UserStatus             string    `json:"userStatus"`
	UserID                 int       `json:"userId"`
	Name                   string    `json:"name"`
	HasLive                bool      `json:"hasLive"`
}

func (t *tokenResp) toToken() *Token {
	return &Token{
		AccessToken:            t.AccessToken,
		ExpirationTime:         t.ExpirationTime,
		PasswordExpirationTime: t.PasswordExpirationTime,
		UserStatus:             t.UserStatus,
		UserID:                 t.UserID,
		Name:                   t.Name,
		HasLive:                t.HasLive,
	}
}

func (r *REST) newToken(ctx context.Context) (*Token, error) {
	buf, err := json.Marshal(r.tokenManager.creds)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		r.baseURL+"/"+accessTokenURL,
		bytes.NewReader(buf),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")

	resp, err := r.h.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, newErrFromResp(resp)
	}

	var t tokenResp
	if err = json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}

	if t.ErrorText != "" {
		return nil, &Err{Status: 200, Body: t.ErrorText}
	}

	newToken := t.toToken()
	r.tokenManager.mu.Lock()
	defer r.tokenManager.mu.Unlock()

	r.tokenManager.token = newToken
	return newToken, nil
}

func (r *REST) refreshToken(ctx context.Context) (*Token, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.baseURL+"/auth/renewAccessToken", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	r.tokenManager.mu.RLock()
	req.Header.Add("authorization", "Bearer "+r.tokenManager.token.AccessToken)
	r.tokenManager.mu.RUnlock()

	resp, err := r.h.Do(req)
	if err != nil {
		return nil, err
	}

	var refreshed tokenResp
	if err = json.NewDecoder(resp.Body).Decode(&refreshed); err != nil {
		return nil, err
	}

	newToken := refreshed.toToken()
	r.tokenManager.mu.Lock()
	defer r.tokenManager.mu.Unlock()

	r.tokenManager.token = newToken
	return newToken, nil
}
