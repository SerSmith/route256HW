package domain

import "context"

func (m *Model) ListOrder(ctx context.Context, orderID int64) (Order, error) {

	order, err := m.DB.GetOrderDetails(ctx, orderID)

	return order, err
}
