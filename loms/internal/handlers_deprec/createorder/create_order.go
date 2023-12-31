package createorder

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
	OrderID int64 `json:"orderID"`
}

type Request struct {
	User  int64              `json:"user"`
	Items []domain.ItemOrder `json:"items"`
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
	logger.Info("create order, request: %+v", req)

	items := make([]domain.ItemOrder, 0, len(req.Items))

	for _, item := range req.Items {
		items = append(items, domain.ItemOrder{
			SKU:   item.SKU,
			Count: item.Count,
		})
	}

	id, err := h.Model.CreateOrder(ctx, req.User, items)
	if err != nil {
		logger.Info("create order: %s", err)
		return Response{}, err
	}
	return Response{
		OrderID: id,
	}, err
}
