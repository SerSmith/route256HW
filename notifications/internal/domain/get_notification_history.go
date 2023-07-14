package domain

import (
	"context"
	"route256/libs/logger"
	"route256/libs/tracer"
	"time"

	"github.com/opentracing/opentracing-go"
)

func (m *Model) GetNotificationHistory(ctx context.Context, req NotificationHistoryRequest) ([]NotificationMem, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "domain/get_notification_history/GetNotificationHistory")
	defer span.Finish()

	Res, ok, err := m.CashDB.Get(ctx, req)

	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, err)
	}

	if !ok {

		logger.Info("Go to Repository")

		Res, err = m.Repository.ReadNotifications(ctx, req)

		if err != nil {
			return nil, tracer.MarkSpanWithError(ctx, err)
		}

		if time.Now().After(req.DateTo) {
			/* Если DateTo в будущем не имеем право его запоминать, так как могут появиться новые сообщения*/
			err = m.CashDB.Set(ctx, req, Res)

			if err != nil {
				return nil, tracer.MarkSpanWithError(ctx, err)
			}

		}

	}

	return Res, err
}
