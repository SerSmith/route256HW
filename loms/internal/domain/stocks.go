package domain

import "context"

func (m *Model) Stocks(ctx context.Context, sku uint32) ([]Stock, error) {
	return []Stock{
		{
			WarehouseID: 1,
			Count:       37,
		},
	}, nil
}
