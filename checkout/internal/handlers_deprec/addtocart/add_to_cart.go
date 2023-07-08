package addtocart

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
}

type Request struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
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
	logger.Info("%+v", req)

	err := h.Model.AddToCart(ctx, req.User, req.SKU, req.Count)
	return Response{}, err
}
