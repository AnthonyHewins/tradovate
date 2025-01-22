package tradovate

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type fanoutMutex struct {
	mu       sync.Mutex
	deadline time.Duration
	acc      int
	channels []*fanout
}

type fanout struct {
	c        chan *rawMsg
	deadline time.Duration
	id       int
}

func (f *fanout) wait(ctx context.Context) (*rawMsg, error) {
	ctx, cancel := context.WithTimeout(ctx, f.deadline)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case v := <-f.c:
			if v == nil {
				return nil, fmt.Errorf("race condition error? received nil *rawMsg")
			}

			return v, nil
		}
	}
}

func (f *fanoutMutex) request() *fanout {
	f.mu.Lock()
	defer f.mu.Unlock()

	c := &fanout{
		c:        make(chan *rawMsg, 1),
		deadline: f.deadline,
		id:       f.acc,
	}
	f.acc++

	f.channels = append(f.channels, c)
	return c
}

func (f *fanoutMutex) pub(r *rawMsg) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, v := range f.channels {
		if v.id != r.ID {
			continue
		}

		v.c <- r
		close(v.c)
		n := len(f.channels)
		f.channels[i] = f.channels[n-1]
		f.channels = f.channels[:n-1]
		return
	}
}
