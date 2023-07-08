package loms

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/loms/loms"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func (c *Client) CreateOrder(ctx context.Context, userID int64, items []domain.ItemOrder) (int64, error) {
	_, cancel := context.WithTimeout(context.Background(), c.wait_time)
	span, ctx := opentracing.StartSpanFromContext(ctx, "main")
	defer func() {
		cancel()
		span.Finish()
	}()

	var Items_for_req []*loms_v1.ItemOrder

	for _, item := range items {
		Items_for_req = append(Items_for_req, &loms_v1.ItemOrder{
			SKU:   item.SKU,
			Count: uint32(item.Count),
		})
	}

	req := &loms_v1.CreateOrderRequest{
		User:  userID,
		Items: Items_for_req,
	}

	resp, err := c.loms.CreateOrder(ctx, req)

	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, errors.Wrap(err, "purchase"))
	}

	return resp.GetOrderID(), nil
}
