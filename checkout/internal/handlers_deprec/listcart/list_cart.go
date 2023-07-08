package listcart

import (
	"context"
	"errors"
	"route256/checkout/internal/domain"
	"route256/libs/logger"
)

type Handler struct {
	Model *domain.Model
}

type ItemCart struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type Response struct {
	Items      []ItemCart `json:"items"`
	TotalPrice uint32     `json:"totalPrice"`
}

type Request struct {
	User int64 `json:"user"`
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrSKUNotFound  = errors.New("SKU not found")
	ErrWrongCount   = errors.New("wrong count")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {

	logger.Info("list cart, request: %+v", req)

	totalPrice, cart, err := h.Model.ListCart(ctx, req.User)
	if err != nil {
		logger.Info("list cart: %s", err)
		return Response{}, err
	}
	resp := Response{
		TotalPrice: totalPrice,
		Items:      make([]ItemCart, 0, len(cart)),
	}

	for _, item := range cart {
		resp.Items = append(resp.Items, ItemCart{
			SKU:   item.SKU,
			Count: item.Count,
			Price: item.Product.Price,
			Name:  item.Product.Name,
		})
	}
	return resp, err
}
