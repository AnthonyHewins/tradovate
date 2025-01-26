package tradovate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func (t *Token) equal(x *Token) bool {
	return (t.AccessToken == x.AccessToken &&
		t.ExpirationTime == x.ExpirationTime &&
		t.PasswordExpirationTime == x.PasswordExpirationTime &&
		t.UserStatus == x.UserStatus &&
		t.UserID == x.UserID &&
		t.Name == x.Name &&
		t.HasLive == x.HasLive)
}

func TestToken(mainTest *testing.T) {
	validTokenResp1 := &tokenResp{
		AccessToken:    "a",
		ExpirationTime: time.Now().Add(time.Hour * 2).UTC(),
	}
	validToken1 := validTokenResp1.toToken()

	tokenDueForRefresh := &Token{
		AccessToken:    "b",
		ExpirationTime: time.Now().Add(time.Hour),
	}
	// expiredToken := Token{AccessToken: "b"}

	testCases := []struct {
		name        string
		start       *Token
		expected    *Token
		expectedErr string

		expectedEndpoint string

		mock       *tokenResp
		mockStatus int
	}{
		{
			name:     "with no access token, fetches new one",
			expected: validToken1,

			expectedEndpoint: accessTokenURL,
			mockStatus:       200,
			mock:             validTokenResp1,
		},
		{
			name:     "with token that needs refresh, it refreshes",
			start:    tokenDueForRefresh,
			expected: validToken1,

			expectedEndpoint: renewTokenURL,
			mockStatus:       200,
			mock:             validTokenResp1,
		},
		{
			name:     "with token that's good for usage, doesnt fetch anything",
			start:    validToken1,
			expected: validToken1,
		},
	}

	for _, tc := range testCases {
		mainTest.Run(tc.name, func(tt *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.expectedEndpoint == "" || !strings.Contains(r.RequestURI, tc.expectedEndpoint) {
					w.WriteHeader(500)
					w.Write([]byte(
						fmt.Sprintf(
							`{"errorText":"failed test: called wrong endpoint %s but wanted %s"}`,
							r.RequestURI,
							tc.expectedEndpoint,
						),
					))
					return
				}

				buf, err := json.Marshal(tc.mock)
				if err != nil {
					panic(err)
				}

				w.WriteHeader(tc.mockStatus)
				w.Write(buf)
			}))
			defer s.Close()

			r := REST{
				baseURL: s.URL,
				tokenManager: tokenManager{
					forceRefreshDeadline: time.Hour,
					token:                tc.start,
				},
				h: &http.Client{},
			}

			actual, actualErr := r.Token(context.Background())
			if tc.expectedErr != "" {
				if actualErr == nil || actualErr.Error() != tc.expectedErr {
					tt.Errorf("wanted error %v but got %v", tc.expectedErr, actualErr)
				}
				return
			}

			if actualErr != nil {
				tt.Errorf("should not have errored but got %v", actualErr)
				return
			}

			if tc.expected.equal(actual) {
				return
			}

			tt.Errorf("invalid token response\nwant: %+v\n got: %+v", tc.expected, actual)
		})
	}
}
