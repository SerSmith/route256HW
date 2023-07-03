package domain

import (
	"context"
	"fmt"
)

func (m *Model) ChangeOrderStatusWithNotification(ctx context.Context, orderID int64, newStatus OrderStatus) error {

	nowStatus, err := m.DB.GetOrderStatus(ctx, orderID)
	if err != nil {
		return fmt.Errorf("GetOrderStatus: %w", err)
	}

	err = m.DB.ChangeOrderStatus(ctx, orderID, newStatus)
	if err != nil {
		return fmt.Errorf("ChangeOrderStatus: %w", err)
	}

	message := StatusChangeMessage{
		OldStatus: string(nowStatus),
		NewStatus: string(newStatus),
		OrderID:   orderID}

	err = m.KP.SendMessage(message)
	if err != nil {
		return fmt.Errorf("sendMessage: %w", err)
	}

	return nil
}
