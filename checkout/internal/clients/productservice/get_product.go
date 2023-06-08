package productservice

import (
	"context"
	"fmt"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/product_service/product_service"
)

func (c *Client) GetProduct(ctx context.Context, sku uint32) (*domain.Product, error) {
	_, cancel := context.WithTimeout(context.Background(), c.wait_time)
	defer cancel()

	response_from_api, err := c.psClient.GetProduct(ctx, &product_service.GetProductRequest{
		Token: c.token,
		Sku:   sku,
	})

	if err != nil {
		return nil, fmt.Errorf("send getProduct request: %w", err)
	}

	response := &domain.Product{Name: response_from_api.GetName(), Price: response_from_api.GetPrice()}

	return response, nil
}
