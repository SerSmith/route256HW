package loms

import (
	"context"
	"log"
	"net/http"
	"route256/checkout/internal/domain"
	"route256/libs/cliwrapper"
)

type StocksRequest struct {
	SKU uint32 `json:"sku"`
}

type StocksResponse struct {
	Stocks []struct {
		WarehouseID int64  `json:"warehouseID"`
		Count       uint64 `json:"count"`
	} `json:"stocks"`
}

func (c *Client) Stocks(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	requestStocks := StocksRequest{SKU: sku}

	ctx, fnCancel := context.WithTimeout(ctx, waittime)
	defer fnCancel()

	responseStocks, err := cliwrapper.RequestAPI[StocksRequest, StocksResponse](ctx, http.MethodGet, c.stockURL, requestStocks)
	if err != nil {
		log.Printf("loms client, get stocks: %s", err)
		log.Printf("stockURL: %s", c.stockURL)
		return nil, err
	}
	result := make([]domain.Stock, 0, len(responseStocks.Stocks))
	for _, v := range responseStocks.Stocks {
		result = append(result, domain.Stock{
			WarehouseID: v.WarehouseID,
			Count:       v.Count,
		})
	}

	return result, nil
}
