package domain

import (
	"context"
	"route256/libs/logger"
	"route256/libs/tracer"
	"route256/libs/workerpool"
	"sync"

	"github.com/opentracing/opentracing-go"
)

const (
	WORKERS_NUM = 5
)

var (
	GetProductUnknownName  = "Unknown"
	GetProductUnknownPrice = uint32(0)
)

func (m *Model) ListCart(ctx context.Context, user int64) (uint32, []ItemCart, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "domain/list_cart/ListCart")
	defer span.Finish()

	OrderItems, err := m.DB.GetCartDB(ctx, user)

	if err != nil {
		return 0, nil, tracer.MarkSpanWithError(ctx, err)
	}

	resChan := make(chan ItemCart, len(OrderItems))

	wp := workerpool.New(WORKERS_NUM)
	wg := sync.WaitGroup{}

	for _, item := range OrderItems {
		wg.Add(1)
		NowItem := item
		err := wp.Run(ctx,
			func(ctx context.Context) {
				defer wg.Done()

				product, err := m.productServiceClient.GetProduct(ctx, NowItem.SKU)

				if err != nil {
					logger.Info("err in runOmeGetProductInstance %w", err)
					product = &Product{Name: GetProductUnknownName,
						Price: GetProductUnknownPrice}
				}

				resChan <- ItemCart{
					SKU:     NowItem.SKU,
					Product: *product,
					Count:   NowItem.Count,
				}
			})
		if err != nil {
			return 0, nil, tracer.MarkSpanWithError(ctx, err)
		}

	}

	wg.Wait()

	var outCart []ItemCart
	var totalPrice uint32

	for range OrderItems {
		oneOutCart := <-resChan

		outCart = append(outCart, oneOutCart)
		totalPrice += oneOutCart.Product.Price * uint32(oneOutCart.Count)
	}

	return totalPrice, outCart, nil
}
