package loms

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/loms/loms"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func (c *Client) Stocks(ctx context.Context, SKU uint32) ([]domain.Stock, error) {
	_, cancel := context.WithTimeout(context.Background(), c.wait_time)
	span, ctx := opentracing.StartSpanFromContext(ctx, "clients/loms")
	defer func() {
		cancel()
		span.Finish()
	}()

	resp, err := c.loms.Stocks(ctx, &loms_v1.StocksRequest{SKU: SKU})
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "stocks"))
	}

	result := make([]domain.Stock, 0, len(resp.Stocks))
	for _, v := range resp.Stocks {
		result = append(result, domain.Stock{
			WarehouseID: v.GetWarehouseID(),
			Count:       v.GetCount(),
		})
	}

	return result, nil
}
