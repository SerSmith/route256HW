package domain

import (
"context"
"fmt"
)

func (m *Model) Purchase(ctx context.Context, userid int64) (int64, error) {

	items, err := m.DB.GetCartDB(ctx, userid)

	if err != nil {
		return 0, fmt.Errorf("err in getCartDB: %v", err)
	}

	if len(items) == 0{
		return 0, fmt.Errorf("Empty cart")
	}


	return m.LomsClient.CreateOrder(ctx, userid, items)
}
