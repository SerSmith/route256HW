package domain

import (
	"context"
	"fmt"
	"route256/libs/logger"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
)

func (m *Model) checkPtoductReservation(ctx context.Context, item ItemOrder) ([]StockInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "domain/create_oeder/checkPtoductReservation")
	defer span.Finish()

	countNotReserved := uint64(item.Count)

	stocksFound, err := m.DB.GetAvailableBySku(ctx, item.SKU)

	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, err)
	}

	var new_reservedStocks []StockInfo

	for _, stock := range stocksFound {

		if countNotReserved <= 0 {
			break
		}

		logger.Info("stock.Count ", stock.Count)
		logger.Info("countNotReserved", countNotReserved)

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
		return nil, tracer.MarkSpanWithError(ctx, fmt.Errorf("Could not reserve enough for sku %v", item.SKU))
	}

	return new_reservedStocks, nil
}

func (m *Model) CreateOrder(ctx context.Context, userID int64, items []ItemOrder) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "domain/create_oeder/CreateOrder")
	defer span.Finish()

	var orderID int64

	err := m.DB.RunRepeatableRead(ctx,
		func(ctxTx context.Context) error {

			OrderID_, err := m.DB.WriteOrder(ctx, items, userID)

			OrderID = OrderID_
			if err != nil {
				return tracer.MarkSpanWithError(ctx, err)
			}

			err = m.ChangeOrderStatusWithNotification(ctx, OrderID, NewStatus)
			if err != nil {
				return tracer.MarkSpanWithError(ctx, err)
			}

			var reservedStocks [][]StockInfo
			for _, item := range items {
				reservedStock, err := m.checkPtoductReservation(ctx, item)
				reservedStocks = append(reservedStocks, reservedStock)

				if err != nil {

					err2 := m.ChangeOrderStatusWithNotification(ctx, OrderID, FailedStatus)

					if err2 != nil {
						return tracer.MarkSpanWithError(ctx, err2)
					}

					return tracer.MarkSpanWithError(ctx, err)
				}

			}

			for _, reservedStock := range reservedStocks {

				err = m.DB.ReserveProducts(ctx, OrderID, reservedStock)

				if err != nil {
					err2 := m.ChangeOrderStatusWithNotification(ctx, OrderID, FailedStatus)
					if err2 != nil {
						return tracer.MarkSpanWithError(ctx, err2)
					}
					return tracer.MarkSpanWithError(ctx, err)
				}

				err = m.DB.MinusAvalibleCount(ctx, reservedStock)

				if err != nil {
					err2 := m.ChangeOrderStatusWithNotification(ctx, OrderID, FailedStatus)
					if err2 != nil {
						return tracer.MarkSpanWithError(ctx, err2)
					}
					return tracer.MarkSpanWithError(ctx, err)
				}
			}

			err = m.ChangeOrderStatusWithNotification(ctx, OrderID, AwaitingPaymentStatus)
			if err != nil {
				return tracer.MarkSpanWithError(ctx, err)
			}

			return nil

		})

	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, err)
	}

	return OrderID, err
}
