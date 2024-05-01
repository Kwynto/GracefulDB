package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	MIN_SIZE_MAP        = 1
	MAX_TIME_CLOSE      = 5
	MICRO_DEFAULT_DELAY = 50
)

type TFnHandler func(ctx context.Context, c *TCloser)

type TCloser struct {
	mu      sync.RWMutex
	funcs   map[string]TFnHandler
	Msgs    []string
	Counter int
}

var StCloseProcs = &TCloser{
	funcs: make(map[string]TFnHandler, MIN_SIZE_MAP),
}

func New() *TCloser {
	return &TCloser{
		funcs: make(map[string]TFnHandler, MIN_SIZE_MAP),
	}
}

func (c *TCloser) AddMsg(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Msgs = append(c.Msgs, fmt.Sprintf("[!] %v", msg))
}

func (c *TCloser) Done() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Counter--
}

func (c *TCloser) AddHandler(f TFnHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs[fmt.Sprint(f)] = f
	c.Counter++
}

func (c *TCloser) DelHandler(f TFnHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := fmt.Sprint(f)
	if _, ok := c.funcs[key]; ok {
		delete(c.funcs, key)
		c.Counter--
	}
}

func (c *TCloser) RunAndDelHandler(f TFnHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	sdCtx, cnl := context.WithTimeout(context.Background(), MAX_TIME_CLOSE*time.Second)
	defer cnl()

	go f(sdCtx, c)

	key := fmt.Sprint(f)
	delete(c.funcs, key)
}

func (c *TCloser) Close(ctx context.Context) error {
	var complete = make(chan struct{}, MIN_SIZE_MAP)

	for _, f := range c.funcs {
		go f(ctx, c)
	}

	go func() {
		time.Sleep(MICRO_DEFAULT_DELAY * time.Millisecond)
		for {
			time.Sleep(MICRO_DEFAULT_DELAY * time.Millisecond)
			if c.Counter <= 0 {
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

	if len(c.Msgs) > 0 {
		return fmt.Errorf("shutdown finished with error(s): %s", strings.Join(c.Msgs, " | "))
	}

	return nil
}

func AddHandler(f TFnHandler) {
	StCloseProcs.AddHandler(f)
}

func DelHandler(f TFnHandler) {
	StCloseProcs.DelHandler(f)
}

func RunAndDelHandler(f TFnHandler) {
	StCloseProcs.RunAndDelHandler(f)
}

func Close(ctx context.Context) error {
	return StCloseProcs.Close(ctx)
}
