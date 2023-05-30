package createorder

import (
	"context"
	"errors"
	"log"
	"route256/loms/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

type OrderItem struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Request struct {
	User  int64       `json:"user"`
	Items []OrderItem `json:"items"`
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrOrderIsEmpty = errors.New("order is empty")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrUserNotFound
	}
	if len(r.Items) == 0 {
		return ErrOrderIsEmpty
	}
	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("create order, request: %+v", req)

	items := make([]domain.ItemOrder, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.ItemOrder{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}
	id, err := h.Model.CreateOrder(ctx, req.User, items)
	if err != nil {
		log.Printf("create order: %s", err)
		return Response{}, err
	}
	return Response{
		OrderID: id,
	}, err
}
