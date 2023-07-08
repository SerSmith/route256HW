package productservice

import (
	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/product_service/product_service"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
)

func (c *Client) GetProduct(ctx context.Context, sku uint32) (*domain.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "clients/productservice/get_product/GetProduct")
	defer span.Finish()

	_, cancel := context.WithTimeout(context.Background(), c.wait_time)
	defer cancel()

	err := c.Limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}

	response_from_api, err := c.psClient.GetProduct(ctx, &product_service.GetProductRequest{
		Token: c.token,
		Sku:   sku,
	})

	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, err)
	}

	response := &domain.Product{Name: response_from_api.GetName(), Price: response_from_api.GetPrice()}

	return response, nil
}
