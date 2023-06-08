package productservice

import (
	"context"
	"log"
	"net/http"
	"route256/checkout/internal/domain"
	"route256/libs/cliwrapper"
)

type GetProductRequest struct {
	Token string `json:"token"`
	SKU   uint32 `json:"sku"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func (c *Client) GetProduct(ctx context.Context, sku uint32) (domain.Product, error) {
	productRequest := GetProductRequest{
		Token: c.token,
		SKU:   sku,
	}

	ctx, fnCancel := context.WithTimeout(ctx, waittime)
	defer fnCancel()

	productResponse, err := cliwrapper.RequestAPI[GetProductRequest, GetProductResponse](ctx, http.MethodPost, c.productPath, productRequest)
	if err != nil {
		log.Printf("product service client, get product: %s", err)
		return domain.Product{}, err
	}

	return domain.Product{
		Name:  productResponse.Name,
		Price: productResponse.Price,
	}, nil
}
