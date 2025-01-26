package tradovate

import (
	"reflect"
	"testing"
)

func TestNewFrame(mainTest *testing.T) {
	testCases := []struct {
		name        string
		arg         []byte
		expected    frame
		expectedErr string
	}{
		{
			name:        "base case",
			expectedErr: ErrEmptyFrame.Error(),
		},
		{
			name: "data msg of a response to request",
			arg:  []byte(`a[{"s":200,"i":23,"d":{"id":65543,"name":"CLZ6","contractMaturityId":6727}}]`),
			expected: dataframe{[]rawMsg{{
				ID:     23,
				Status: 200,
				Data:   []byte(`{"id":65543,"name":"CLZ6","contractMaturityId":6727}`),
			}}},
		},
		{
			name: "data msg of event publish",
			arg:  []byte(`a[{"e":"props","d":{"entityType":"order","eventType":"Created","entity":{"id":210518,"accountId":25,"contractId":560901,"timestamp":"2016-11-04T00:02:36.626Z","action":"Sell","ordStatus":"PendingNew","admin":false}}}]`),
			expected: dataframe{[]rawMsg{{
				Event: frameEventProps,
				Data:  []byte(`{"entityType":"order","eventType":"Created","entity":{"id":210518,"accountId":25,"contractId":560901,"timestamp":"2016-11-04T00:02:36.626Z","action":"Sell","ordStatus":"PendingNew","admin":false}}`),
			}}},
		},
	}

	for _, tc := range testCases {
		mainTest.Run(tc.name, func(tt *testing.T) {
			actual, actualErr := newFrame(tc.arg)

			if tc.expectedErr != "" {
				if actualErr == nil || actualErr.Error() != tc.expectedErr {
					tt.Errorf("wanted err %v but got %v", tc.expectedErr, actualErr)
				}
				return
			}

			if actualErr != nil {
				tt.Errorf("wanted no error, but got %v", actualErr)
				return
			}

			if !reflect.DeepEqual(tc.expected, actual) {
				tt.Errorf("want: %v\n got: %v", tc.expected, actual)
			}
		})
	}
}
