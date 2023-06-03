package loms

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/loms/loms"
)

func (c *Client) Stocks(ctx context.Context, SKU uint32) ([]domain.Stock, error) {
	_, cancel := context.WithTimeout(context.Background(), c.wait_time)
	defer cancel()

	fmt.Println("I am here loms client 2.3")

	resp, err := c.loms.Stocks(ctx, &loms_v1.StocksRequest{SKU: SKU})
	if err != nil {
		return nil, fmt.Errorf("stocks: %w", err)
	}

	result := make([]domain.Stock, 0, len(resp.Stocks))
	for _, v := range resp.Stocks {
		result = append(result, domain.Stock{
			WarehouseID: v.GetWarehouseID(),
			Count:       v.GetCount(),
		})
	}

	return result, nil
}
