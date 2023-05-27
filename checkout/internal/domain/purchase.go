package domain

import "context"

func (m *Model) Purchase(ctx context.Context, userid int64) (int64, error) {

	items := []ItemOrder{
		{
			SKU:   1148162,
			Count: 2,
		},
		{
			SKU:   32885918,
			Count: 1,
		},
	}

	return m.loms.CreateOrder(ctx, userid, items)
}
