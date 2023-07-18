package domain

import (
	"context"
	"fmt"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
)

func (m *Model) Purchase(ctx context.Context, userid int64) (int64, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "domain/purchase/Purchase")
	defer span.Finish()

	items, err := m.DB.GetCartDB(ctx, userid)

	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, err)
	}

	if len(items) == 0 {
		return 0, tracer.MarkSpanWithError(ctx, fmt.Errorf("Empty cart"))
	}

	orderID, err := m.LomsClient.CreateOrder(ctx, userid, items)

	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, err)
	}

	err = m.DB.WipeCartDB(ctx, userid)

	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, err)
	}

	return orderID, nil
}
