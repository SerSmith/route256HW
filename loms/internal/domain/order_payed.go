package domain

import (
	"context"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
)

func (m *Model) OrderPayed(ctx context.Context, orderID int64) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "domain/order_payed/OrderPayed")
	defer span.Finish()

	stocks, err := m.DB.GetReservedByOrderID(ctx, orderID)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	err = m.DB.BuyProducts(ctx, stocks)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	err = m.DB.UnreserveProducts(ctx, orderID)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	err = m.ChangeOrderStatusWithNotification(ctx, orderID, PayedStatus)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil
}
