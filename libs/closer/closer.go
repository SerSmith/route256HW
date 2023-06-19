package closer

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Closer struct {
	mu    sync.Mutex
	funcs []CloseFunc
}

type CloseFunc func(ctx context.Context) error

func (c *Closer) Add(f CloseFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	done := make(chan struct{}, 1)
	go func() {
		for _, f := range c.funcs {
			if err := f(ctx); err != nil {
				log.Println(err)
			}
		}

		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	}
}
