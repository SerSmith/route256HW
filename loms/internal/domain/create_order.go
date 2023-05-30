package domain

import "context"

func (m *Model) CreateOrder(ctx context.Context, userID int64, items []ItemOrder) (int64, error) {
	return 1, nil
}
