package deletefromcart

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/domain"
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
	ErrSKUNotFound  = errors.New("SKU not found")
	ErrWrongCount   = errors.New("wrong count")
)

func (r Request) Validate() error {
	if r.User == 0 {
		return ErrUserNotFound
	}

	if r.SKU == 0 {
		return ErrSKUNotFound
	}

	return nil
}

func (h *Handler) Handle(ctx context.Context, req Request) (Response, error) {
	log.Printf("delete from cart, request: %+v", req)

	err := h.Model.DeleteFromCart(ctx, req.User, req.SKU, req.Count)

	resp := Response{}
	return resp, err
}
