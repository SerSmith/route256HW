package domain

import (
	"context"
	"fmt"
)

func (m *Model) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {

	countCart, err := m.DB.GetCartQauntDB(ctx, user, sku)

	if err != nil {
		return fmt.Errorf("getCartDB: %w", err)
	}

	if countCart < count {
		return fmt.Errorf("Trying to delete from cart more then cart have")
	}

	err = m.DB.DeleteFromCartDB(ctx, user, sku, count)

	if err != nil {
		return fmt.Errorf("err in deleteFromCartDB %w", err)
	}

	return nil
}
