package purchase

import (
	"context"
	"errors"
	"route256/checkout/internal/domain"
	"route256/libs/logger"
)

type Handler struct {
	Model *domain.Model
}

type Response struct {
	OrderID int64 `json:"orderID"`
}

type Request struct {
	User int64 `json:"user"`
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	logger.Info("purchase, request: %+v", req)

	orderID, err := h.Model.Purchase(ctx, req.User)
	if err != nil {
		logger.Info("purchase: %s", err)
		return Response{}, err
	}

	resp := Response{OrderID: orderID}

	return resp, err
}
