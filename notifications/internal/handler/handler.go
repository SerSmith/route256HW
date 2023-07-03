package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/notifications/internal/domain"

	"github.com/Shopify/sarama"
)

type Handler struct {
	model *domain.Model
}

func NewHandler(model *domain.Model) *Handler {
	return &Handler{model: model}
}

func (h *Handler) Notify(ctx context.Context, message *sarama.ConsumerMessage) error {

	scm := domain.StatusChangeMessage{}
	err := json.Unmarshal(message.Value, &scm)
	if err != nil {
		return fmt.Errorf("Notify json.Unmarshal: %w", err)
	}

	text := fmt.Sprintf("Order %d, has changed status from %s to %s", scm.OrderID, string(scm.OldStatus), string(scm.NewStatus))

	err = h.model.Messenger.SendMessage(ctx, text)

	if err != nil {
		return fmt.Errorf("Notify r.model.Messenger.SendMessage: %w", err)
	}

	return nil
}
