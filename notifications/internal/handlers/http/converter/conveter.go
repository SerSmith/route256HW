package converter

import (
	"context"
	"route256/libs/tracer"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/handlers/http/schema"
	"time"
)

func ReqD2S(ctx context.Context, reqSchema []domain.NotificationMem) (schema.Response, error) {

	out := make([]schema.ResponseEl, 0, len(reqSchema))

	for _, rs := range reqSchema {

		resp := schema.ResponseEl{OldStatus: rs.ChangeStatus.OldStatus,
			NewStatus: rs.ChangeStatus.NewStatus,
			OrderID:   rs.ChangeStatus.OrderID,
			UserID:    rs.ChangeStatus.UserID,
			DT:        rs.DT}

		out = append(out, resp)
	}

	return schema.Response{Data: out}, nil
}

func ReqS2D(ctx context.Context, req schema.Request) (domain.NotificationHistoryRequest, error) {

	DateFrom, err := time.Parse("2006-01-02", req.DateFrom)

	if err != nil {
		return domain.NotificationHistoryRequest{}, tracer.MarkSpanWithError(ctx, err)
	}

	DateTo, err := time.Parse("2006-01-02", req.DateTo)

	if err != nil {
		return domain.NotificationHistoryRequest{}, tracer.MarkSpanWithError(ctx, err)
	}

	return domain.NotificationHistoryRequest{UserID: req.UserID,
		DateFrom: DateFrom,
		DateTo:   DateTo}, nil
}
