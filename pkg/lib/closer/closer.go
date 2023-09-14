package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Handler func(ctx context.Context, c *Closer)

type Closer struct {
	mu      sync.RWMutex
	funcs   []Handler
	msgs    []string
	counter int
}

func (c *Closer) AddMsg(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.msgs = append(c.msgs, fmt.Sprintf("[!] %v", msg))
}

func (c *Closer) Done() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.counter--
}

func (c *Closer) AddHandler(f Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
	c.counter++
}

func (c *Closer) Close(ctx context.Context) error {
	var complete = make(chan struct{}, 1)

	for _, f := range c.funcs {
		go f(ctx, c)
	}

	go func() {
		time.Sleep(250 * time.Millisecond)
		for {
			time.Sleep(50 * time.Millisecond)
			if c.counter <= 0 {
				complete <- struct{}{}
				break
			}
		}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.msgs) > 0 {
		return fmt.Errorf("shutdown finished with error(s): %s", strings.Join(c.msgs, " | "))
	}

	return nil
}
