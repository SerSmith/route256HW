package orderpayed

import (
	"context"
	"errors"
	"route256/libs/logger"
	"route256/loms/internal/domain"
)

type Handler struct {
	Model *domain.Model
}

type Response struct{}

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
	logger.Info("order payed, request: %+v", req)
	err := h.Model.OrderPayed(ctx, req.OrderID)
	return Response{}, err
}
