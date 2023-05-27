package loms

import (
	"context"
	"log"
	"net/http"
	"route256/checkout/internal/domain"
	"route256/libs/cliwrapper"
)

type ItemOrder struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type RequestCreateOrder struct {
	User  int64       `json:"user"`
	Items []ItemOrder `json:"items"`
}

type ResponseCreateOrder struct {
	OrderID int64 `json:"orderID"`
}

func (c *Client) CreateOrder(ctx context.Context, user int64, items []domain.ItemOrder) (int64, error) {
	var (
		requestCreateOrder = RequestCreateOrder{
			User: user,
		}
		reqItems = make([]ItemOrder, 0, len(items))
	)

	for _, item := range items {
		reqItems = append(reqItems, ItemOrder{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}
	requestCreateOrder.Items = reqItems

	ctx, fnCancel := context.WithTimeout(ctx, waittime)
	defer fnCancel()

	responseCreateOrder, err := cliwrapper.RequestAPI[RequestCreateOrder, ResponseCreateOrder](ctx, http.MethodGet, c.createOrderURL, requestCreateOrder)
	if err != nil {
		log.Printf("loms client, create order: %s", err)
		log.Printf("createOrderPath: %s", c.createOrderURL)
		return 0, err
	}

	return responseCreateOrder.OrderID, nil
}
