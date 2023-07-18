package domain

import (
	"context"
	"fmt"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
)

func (m *Model) CancelOrder(ctx context.Context, orderID int64) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "closer/Close")
	defer span.Finish()

	err := m.DB.RunRepeatableRead(ctx,
		func(ctxTx context.Context) error {

			status, err := m.DB.GetOrderStatus(ctx, orderID)

			if err != nil {
				return fmt.Errorf("GetOrderStatus: %s", err)
			}

			if status != AwaitingPaymentStatus {
				return fmt.Errorf("Wrong order status")
			}

			stocks, err := m.DB.GetReservedByOrderID(ctx, orderID)
			if err != nil {
				return fmt.Errorf("GetReservedByOrderID: %s", err)
			}

			err = m.DB.PlusAvalibleCount(ctx, stocks)
			if err != nil {
				return fmt.Errorf("BuyProducts: %s", err)
			}

			err = m.DB.UnreserveProducts(ctx, orderID)
			if err != nil {
				return fmt.Errorf("UnreserveProducts: %s", err)
			}

			err = m.ChangeOrderStatusWithNotification(ctx, orderID, CanceledStatus)
			if err != nil {
				return fmt.Errorf("ChangeOrderStatusWithNotification: %s", err)
			}

			return nil
		})

	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil
}
