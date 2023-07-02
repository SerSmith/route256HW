package loms

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/loms/loms"
)

func (c *Client) CreateOrder(ctx context.Context, userID int64, items []domain.ItemOrder) (int64, error) {
	_, cancel := context.WithTimeout(context.Background(), c.wait_time)
	defer cancel()

	var Items_for_req []*loms_v1.ItemOrder

	for _, item := range items {
		Items_for_req = append(Items_for_req, &loms_v1.ItemOrder{
			SKU:   item.SKU,
			Count: uint32(item.Count),
		})
	}

	req := &loms_v1.CreateOrderRequest{
		User:  userID,
		Items: Items_for_req,
	}

	resp, err := c.loms.CreateOrder(ctx, req)

	if err != nil {
		return 0, fmt.Errorf("purchase: %w", err)
	}

	return resp.GetOrderID(), nil
}
