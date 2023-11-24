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
	funcs   map[string]Handler
	msgs    []string
	counter int
}

var CloseProcs = &Closer{
	funcs: make(map[string]Handler, 1),
}

func New() *Closer {
	return &Closer{
		funcs: make(map[string]Handler, 1),
	}
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

	// c.funcs = append(c.funcs, f)
	c.funcs[fmt.Sprint(f)] = f
	c.counter++
}

func (c *Closer) DelHandler(f Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprint(f)
	if _, ok := c.funcs[key]; ok {
		delete(c.funcs, key)
		c.counter--
	}
}

func (c *Closer) RunAndDelHandler(f Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	sdCtx, cnl := context.WithTimeout(context.Background(), 5*time.Second)
	defer cnl()

	go f(sdCtx, c)

	key := fmt.Sprint(f)
	if _, ok := c.funcs[key]; ok {
		delete(c.funcs, key)
		c.counter--
	}
}

func (c *Closer) Close(ctx context.Context) error {
	var complete = make(chan struct{}, 1)

	for _, f := range c.funcs {
		go f(ctx, c)
	}

	go func() {
		time.Sleep(50 * time.Millisecond)
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

func AddHandler(f Handler) {
	CloseProcs.AddHandler(f)
}

func DelHandler(f Handler) {
	CloseProcs.DelHandler(f)
}

func RunAndDelHandler(f Handler) {
	CloseProcs.RunAndDelHandler(f)
}

func Close(ctx context.Context) error {
	return CloseProcs.Close(ctx)
}
