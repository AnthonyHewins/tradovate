package tradovate

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type fanoutMutex struct {
	mu       sync.Mutex
	timeout  time.Duration
	acc      int
	channels []*socketReq
}

type socketReq struct {
	c        chan *rawMsg
	deadline time.Time
	id       int
}

func (f *socketReq) wait(userCtx, connCtx context.Context) (*rawMsg, error) {
	ctx, cancel := context.WithDeadline(userCtx, f.deadline)
	defer cancel()

	for {
		select {
		case <-connCtx.Done():
			return nil, fmt.Errorf("websocket connection killed: %w", ctx.Err())
		case <-userCtx.Done():
			return nil, ctx.Err()
		case v := <-f.c:
			if v == nil {
				return nil, fmt.Errorf("channel closed early; timeout")
			}

			return v, nil
		}
	}
}

func (f *fanoutMutex) request() *socketReq {
	f.mu.Lock()
	defer f.mu.Unlock()

	c := &socketReq{
		c:        make(chan *rawMsg, 1),
		deadline: time.Now().Add(f.timeout),
		id:       f.acc,
	}
	f.acc++

	f.channels = append(f.channels, c)
	return c
}

func (f *fanoutMutex) pub(r *rawMsg) {
	f.mu.Lock()
	defer f.mu.Unlock()

	goodPtr, n := 0, len(f.channels)
	for i := 0; i < n; i++ {
		v := f.channels[i]
		if v.deadline.Before(time.Now()) {
			close(v.c)
			continue
		}

		if v.id == r.ID {
			v.c <- r
			close(v.c)
			continue
		}

		f.channels[goodPtr] = v
		goodPtr++
	}

	f.channels = f.channels[:goodPtr]
}
