package closer

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"sync"

	"route256/libs/tracer"
	"github.com/opentracing/opentracing-go"
)

type Closer struct {
	mu	sync.Mutex
	funcs []CloseFunc
}

type CloseFunc func(ctx context.Context) error

func (c *Closer) Add(f CloseFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "closer/Close")
	defer span.Finish()

	c.mu.Lock()
	defer c.mu.Unlock()

	done := make(chan struct{}, 1)
	go func() {
		for _, f := range c.funcs {
			if err := f(ctx); err != nil {
				logger.Error(err)
			}
		}

		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("shutdown cancelled: %v", ctx.Err()))
	}
}
