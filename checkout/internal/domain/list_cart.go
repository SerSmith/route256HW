package domain

import (
	"context"
)

func (m *Model) ListCart(ctx context.Context, user int64) (uint32, []ItemCart, error) {
	outCart := []ItemCart{
		{
			SKU:   1076963,
			Count: 1,
		},
		{
			SKU:   1148162,
			Count: 2,
		},
		{
			SKU:   1625903,
			Count: 3,
		},
	}

	var totalPrice uint32

	for i, item := range outCart {
		product, err := m.productServiceClient.GetProduct(ctx, item.SKU)
		if err != nil {
			return 0, nil, err
		}

		outCart[i].Product = product
		totalPrice += product.Price * uint32(item.Count)
	}

	return totalPrice, outCart, nil
}
