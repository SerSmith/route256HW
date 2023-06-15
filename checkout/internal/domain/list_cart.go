package domain

import (
	"context"
	"fmt"
	"log"
)

func (m *Model) ListCart(ctx context.Context, user int64) (uint32, []ItemCart, error) {
	OrderItems, err := m.DB.GetCartDB(ctx, user)

	if err != nil{
		return 0, nil, fmt.Errorf("err in GetCartDB: %v", err)
	}

	var totalPrice uint32

	var outCart []ItemCart

	for _, item := range OrderItems {
		product, err := m.productServiceClient.GetProduct(ctx, item.SKU)

		if err != nil {
			log.Print("err in productServiceClient.GetProduct ", err)
			product = &Product{Name:	"Unknown",
							  Price:	0}
		}

		oneOutCart :=   ItemCart{SKU:		item.SKU,
								Count:		item.Count,
								Product:	*product}

		
		outCart = append(outCart, oneOutCart)
		totalPrice += product.Price * uint32(item.Count)
	}

	return totalPrice, outCart, nil
}
