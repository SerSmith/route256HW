package listorder

import (
	"context"
	"errors"
	"route256/libs/logger"
	"route256/loms/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	User   int64       `json:"user"`
	Items  []ItemOrder `json:"items"`
	Status string      `json:"status"`
}

type ItemOrder struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Request struct {
	OrderID int64 `json:"orderID"`
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

func (r Request) Validate() error {
	if r.OrderID == 0 {
		return ErrOrderNotFound
	}
	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	logger.Info("list order, request: %+v", req)
	order, err := h.Model.ListOrder(ctx, req.OrderID)
	if err != nil {
		logger.Info("list order: %s", err)
		return Response{}, err
	}
	items := make([]ItemOrder, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, ItemOrder{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}
	return Response{
		User:   order.User,
		Items:  items,
		Status: string(order.Status),
	}, nil
}
