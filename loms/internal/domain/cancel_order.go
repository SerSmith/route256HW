package domain

import (
	"context"
	"fmt"
	"log"
)

func (m *Model) CancelOrder(ctx context.Context, orderID int64) error {

	err := m.DB.RunRepeatableRead(ctx,
		func(ctxTx context.Context) error {

			status, err := m.DB.GetOrderStatus(ctx, orderID)

			if err != nil {
				log.Fatalf("GetOrderStatus: %s", err)
			}

			if status != AwaitingPaymentStatus {
				log.Fatalf("Wrong order status")
			}

			stocks, err := m.DB.GetReservedByOrderID(ctx, orderID)
			if err != nil {
				log.Fatalf("GetReservedByOrderID: %s", err)
			}

			err = m.DB.PlusAvalibleCount(ctx, stocks)
			if err != nil {
				log.Fatalf("BuyProducts: %s", err)
			}

			err = m.DB.UnreserveProducts(ctx, orderID)
			if err != nil {
				log.Fatalf("UnreserveProducts: %s", err)
			}

			err = m.ChangeOrderStatusWithNotification(ctx, orderID, CanceledStatus)
			if err != nil {
				log.Fatalf("ChangeOrderStatusWithNotification: %s", err)
			}

			return nil
		})

	if err != nil {
		return fmt.Errorf("err in RunRepeatableRead: %w", err)
	}

	return nil
}
