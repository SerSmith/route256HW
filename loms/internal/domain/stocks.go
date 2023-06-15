package domain

import (
"context"
"fmt"
)

func (m *Model) Stocks(ctx context.Context, sku uint32) ([]Stock, error) {

	stocks, err := m.DB.GetAvailableBySku(ctx, sku)

	if err != nil{
		return nil, fmt.Errorf("GetAvailable: %s", err)
	}

	return stocks, nil
}
