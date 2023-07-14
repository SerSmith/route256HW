package grpcHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/libs/tracer"
	"route256/notifications/internal/domain"
	"time"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
)

type Handler struct {
	model *domain.Model
}

func NewHandler(model *domain.Model) *Handler {
	return &Handler{model: model}
}

func (h *Handler) Notify(ctx context.Context, message *sarama.ConsumerMessage) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "handler/Notify")
	defer span.Finish()

	scm := domain.StatusChangeMessage{}
	err := json.Unmarshal(message.Value, &scm)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	text := fmt.Sprintf("Order %d, has changed status from %s to %s", scm.OrderID, string(scm.OldStatus), string(scm.NewStatus))

	err = h.model.Messenger.SendMessage(ctx, text)

	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	err = h.model.Repository.WriteNotification(ctx, scm, time.Now())

	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil
}
