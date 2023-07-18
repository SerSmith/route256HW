package domain

import (
	"context"
	"errors"
	"route256/libs/logger"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
)

const ServiceName = "Checkout"

var (
	ErrStockInsufficient = errors.New("No such SKUs on stocks")
)

func (m *Model) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "domain/add_ro_cart/get_product/AddToCart")
	defer span.Finish()

	stocks, err := m.LomsClient.Stocks(ctx, sku)

	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	logger.Info("stocks: %v", stocks)

	ok := false

	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			ok = true
			break
		}
	}

	if ok {
		err = m.DB.AddToCartDB(ctx, user, sku, count)

		if err != nil {
			return tracer.MarkSpanWithError(ctx, err)
		}

	} else {
		return ErrStockInsufficient
	}

	return nil

}
