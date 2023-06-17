package workerpool

import (
	"context"
)

type WorkerPool chan struct{}

func New(limit int) WorkerPool {
	return make(chan struct{}, limit)
}

func (wp WorkerPool) Run(ctx context.Context, f func(ctx context.Context)) error {
	select {
	case <-ctx.Done():
	case wp <- struct{}{}:
		go func() {
			f(ctx)
			select {
				case <-ctx.Done():
				case <-wp:
			}
		}()
	}
	return ctx.Err()
}