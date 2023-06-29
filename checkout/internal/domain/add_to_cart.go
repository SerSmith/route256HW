package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
)

var (
	ErrStockInsufficient = errors.New("No such SKUs on stocks")
)

func (m *Model) AddToCart(ctx context.Context, user int64, sku uint32, count uint16) error {

	stocks, err := m.LomsClient.Stocks(ctx, sku)

	if err != nil {
		return fmt.Errorf("get stocks: %w", err)
	}

	log.Printf("stocks: %v", stocks)

	ok := false

	counter := int64(count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			ok = true
			break
		}
	}

	if ok {
		err = m.DB.AddToCartDB(ctx, user, sku, count)

		if err != nil {
			return fmt.Errorf("AddToCartDB: %w", err)
		}

	} else {
		return ErrStockInsufficient
	}

	return nil

}
