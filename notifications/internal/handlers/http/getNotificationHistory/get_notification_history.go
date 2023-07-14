package getNotificationHistory

import (
	"context"
	"route256/libs/logger"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/handlers/http/converter"
	"route256/notifications/internal/handlers/http/schema"
)

type Handler struct {
	Model *domain.Model
}

func (h *Handler) Handle(ctx context.Context, req schema.Request) (schema.Response, error) {
	logger.Info("%+v", req)

	reqConverted, err := converter.ReqS2D(ctx, req)

	if err != nil {
		return schema.Response{}, err
	}

	res, err := h.Model.GetNotificationHistory(ctx, reqConverted)

	if err != nil {
		return schema.Response{}, err
	}

	resSchema, err := converter.ReqD2S(ctx, res)

	if err != nil {
		return schema.Response{}, err
	}

	return resSchema, nil
}
