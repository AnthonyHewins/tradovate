package tradovate

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func (f *fanoutMutex) equal(x *fanoutMutex) bool {
	if len(f.channels) != len(x.channels) {
		return false
	}

	for i := 0; i < len(f.channels); i++ {
		if eq := f.channels[i].equal(x.channels[i]); !eq {
			return false
		}
	}

	return f.acc == x.acc && f.deadline == x.deadline
}

func (f *fanoutMutex) String() string {
	var sb strings.Builder

	for i, v := range f.channels {
		sb.WriteString(v.String())
		if i != len(f.channels)-1 {
			sb.WriteRune(',')
		}
	}

	return fmt.Sprintf(
		"{Acc:%d,\n\tChannels:[%s],\n\tDeadline: %s}",
		f.acc, sb.String(), f.deadline,
	)
}

func (s *socketReq) equal(x *socketReq) bool {
	return (x.deadline.Truncate(time.Millisecond) == s.deadline.Truncate(time.Millisecond) &&
		x.id == s.id)
}

func (s *socketReq) String() string {
	return fmt.Sprintf(
		"{Deadline: %s, ID: %d}",
		s.deadline, s.id,
	)
}

var futureDeadline = time.Now().Add(time.Second).Truncate(time.Second)

func newValidReq(id int) *socketReq {
	return &socketReq{id: id, deadline: futureDeadline, c: make(chan *rawMsg, 1)}
}

func newInvalidReq(id int) *socketReq {
	return &socketReq{id: id, c: make(chan *rawMsg, 1)}
}

func TestPub(mainTest *testing.T) {

	testCases := []struct {
		name  string
		msg   *rawMsg
		start *fanoutMutex
		end   *fanoutMutex
	}{
		{
			name:  "base case",
			msg:   &rawMsg{},
			start: &fanoutMutex{},
			end:   &fanoutMutex{},
		},
		{
			name: "removes an event that's past its deadline in case of a write error",
			msg:  &rawMsg{},
			start: &fanoutMutex{
				channels: []*socketReq{{c: make(chan *rawMsg)}},
			},
			end: &fanoutMutex{},
		},
		{
			name: "removes a channel after delivering its msg",
			msg:  &rawMsg{},
			start: &fanoutMutex{
				channels: []*socketReq{
					{deadline: time.Now().Add(time.Second), c: make(chan *rawMsg, 1)},
				},
			},
			end: &fanoutMutex{},
		},
		{
			name: "stable in removal of elements",
			msg:  &rawMsg{ID: 4},
			start: &fanoutMutex{
				channels: []*socketReq{
					newValidReq(1),
					newValidReq(2),
					newValidReq(3),
					newValidReq(4),
					newValidReq(5),
					newValidReq(6),
				},
			},
			end: &fanoutMutex{
				channels: []*socketReq{
					newValidReq(1),
					newValidReq(2),
					newValidReq(3),
					newValidReq(5),
					newValidReq(6),
				},
			},
		},
		{
			name: "example with various removals (unrealistic state, but also still recovers)",
			msg:  &rawMsg{ID: 27},
			start: &fanoutMutex{
				channels: []*socketReq{
					newValidReq(3),
					newValidReq(1),
					newInvalidReq(93452),
					newValidReq(2),
					newInvalidReq(56592),
					newValidReq(5),
					newInvalidReq(92234),
					newInvalidReq(9432),
					newInvalidReq(9122),
					newInvalidReq(92),
					newValidReq(26),
					newInvalidReq(94332),
					newInvalidReq(11111192),
					newValidReq(21),
					newInvalidReq(912232),
					newInvalidReq(9122),
					newValidReq(27),
					newInvalidReq(9562),
					newInvalidReq(911112),
					newValidReq(24),
					newValidReq(6),
					newValidReq(21),
				},
			},
			end: &fanoutMutex{
				channels: []*socketReq{
					newValidReq(3),
					newValidReq(1),
					newValidReq(2),
					newValidReq(5),
					newValidReq(26),
					newValidReq(21),
					newValidReq(24),
					newValidReq(6),
					newValidReq(21),
				},
			},
		},
	}

	for _, tc := range testCases {
		mainTest.Run(tc.name, func(tt *testing.T) {
			tc.start.pub(tc.msg)
			if tc.start.equal(tc.end) {
				return
			}

			tt.Errorf("final state not correct\nexpected: %+v\nactual: %+v", tc.end, tc.start)
		})
	}
}
