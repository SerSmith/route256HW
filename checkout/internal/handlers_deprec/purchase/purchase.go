package purchase

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
	log.Printf("purchase, request: %+v", req)

	orderID, err := h.Model.Purchase(ctx, req.User)
	if err != nil {
		log.Printf("purchase: %s", err)
		return Response{}, err
	}

	resp := Response{OrderID: orderID}

	return resp, err
}
