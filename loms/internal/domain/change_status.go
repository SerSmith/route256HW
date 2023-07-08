package domain

import (
	"context"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
)

func (m *Model) ChangeOrderStatusWithNotification(ctx context.Context, orderID int64, newStatus OrderStatus) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "domain/change_status/ChangeOrderStatusWithNotification")
	defer span.Finish()

	nowStatus, err := m.DB.GetOrderStatus(ctx, orderID)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	err = m.DB.ChangeOrderStatus(ctx, orderID, newStatus)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	message := StatusChangeMessage{
		OldStatus: string(nowStatus),
		NewStatus: string(newStatus),
		OrderID:   orderID}

	err = m.KP.SendMessage(message)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil
}
