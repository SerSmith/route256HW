package domain

import (
	"context"
	"log"
)

func (m *Model) OrderPayed(ctx context.Context, orderID int64) error {

	stocks, err := m.DB.GetReservedByOrderID(ctx, orderID)
	if err != nil {
		log.Fatalf("GetReservedByOrderID: %s", err)
	}

	err = m.DB.BuyProducts(ctx, stocks)
	if err != nil {
		log.Fatalf("BuyProducts: %s", err)
	}

	err = m.DB.UnreserveProducts(ctx, orderID)
	if err != nil {
		log.Fatalf("UnreserveProducts: %s", err)
	}

	err = m.ChangeOrderStatusWithNotification(ctx, orderID, PayedStatus)
	if err != nil {
		log.Fatalf("ChangeOrderStatus: %s", err)
	}

	return nil
}
