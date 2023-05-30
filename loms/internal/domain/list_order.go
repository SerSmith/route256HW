package domain

import "context"

func (m *Model) ListOrder(ctx context.Context, orderID int64) (Order, error) {
	return Order{
		User: 1,
		Items: []*ItemOrder{
			{
				SKU:   5423,
				Count: 43,
			},
		},
		Status: NewStatus,
	}, nil
}
