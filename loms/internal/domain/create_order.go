package domain

import (
	"context"
	"fmt"
	"log"
)

func (m *Model) checkPtoductReservation(ctx context.Context, item ItemOrder) ([]StockInfo, error) {
	countNotReserved := uint64(item.Count)

	stocksFound, err := m.DB.GetAvailableBySku(ctx, item.SKU)

	if err != nil {
		return nil, fmt.Errorf("err in GetAvailableBySku: %v", err)
	}

	var new_reservedStocks []StockInfo

	for _, stock := range stocksFound {

		if countNotReserved <= 0 {
			break
		}

		log.Println("stock.Count ", stock.Count)
		log.Println("countNotReserved", countNotReserved)

		var countReserved uint64
		if stock.Count >= countNotReserved {
			countReserved = countNotReserved
		} else {
			countReserved = uint64(stock.Count)
		}

		countNotReserved -= countReserved

		new_reservedStocks = append(new_reservedStocks,
			StockInfo{SKU: int64(item.SKU),
				WarehouseID: stock.WarehouseID,
				Count:       countReserved})

	}

	/* Если не смогли найти достаточно товаров */
	if countNotReserved > 0 {
		return nil, fmt.Errorf("Could not reserve enough for sku %d", item.SKU)
	}

	return new_reservedStocks, nil
}

func (m *Model) CreateOrder(ctx context.Context, userID int64, items []ItemOrder) (int64, error) {

	var OrderID int64

	err := m.DB.RunRepeatableRead(ctx,
		func(ctxTx context.Context) error {

			OrderID_, err := m.DB.WriteOrder(ctx, items, userID)

			OrderID = OrderID_
			if err != nil {
				return fmt.Errorf("err in WriteOrderItems err %v", err)
			}

			err = m.ChangeOrderStatusWithNotification(ctx, OrderID, NewStatus)
			if err != nil {
				return fmt.Errorf("err in ChangeOrderStatus err %v", err)
			}

			var reservedStocks [][]StockInfo
			for _, item := range items {
				reservedStock, err := m.checkPtoductReservation(ctx, item)
				reservedStocks = append(reservedStocks, reservedStock)

				if err != nil {

					err2 := m.ChangeOrderStatusWithNotification(ctx, OrderID, FailedStatus)

					if err2 != nil {
						return fmt.Errorf("err in ChangeOrderStatus err %v", err2)
					}

					return fmt.Errorf("err in checkPtoductReservation err %v", err)
				}

			}

			for _, reservedStock := range reservedStocks {

				err = m.DB.ReserveProducts(ctx, OrderID, reservedStock)

				if err != nil {
					err2 := m.ChangeOrderStatusWithNotification(ctx, OrderID, FailedStatus)
					if err2 != nil {
						return fmt.Errorf("err in ChangeOrderStatus err %v", err2)
					}
					return fmt.Errorf("err in ReserveProducts err %v", err)
				}

				err = m.DB.MinusAvalibleCount(ctx, reservedStock)

				if err != nil {
					err2 := m.ChangeOrderStatusWithNotification(ctx, OrderID, FailedStatus)
					if err2 != nil {
						return fmt.Errorf("err in ChangeOrderStatus err %v", err2)
					}
					return fmt.Errorf("err in MinusAvalibleCount err %v", err)
				}
			}

			err = m.ChangeOrderStatusWithNotification(ctx, OrderID, AwaitingPaymentStatus)
			if err != nil {
				return fmt.Errorf("err in ChangeOrderStatus err %v", err)
			}

			return nil

		})

	if err != nil {
		return 0, fmt.Errorf("err in RunRepeatableRead: %w", err)
	}

	return OrderID, err
}
