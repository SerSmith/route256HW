package domain

import (
	"context"
	"fmt"
	"route256/libs/tracer"

	"github.com/opentracing/opentracing-go"
)

func (m *Model) DeleteFromCart(ctx context.Context, user int64, sku uint32, count uint16) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "domain/delete_from_cart/DeleteFromCart")
	defer span.Finish()

	countCart, err := m.DB.GetCartQauntDB(ctx, user, sku)

	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	if countCart < count {
		return tracer.MarkSpanWithError(ctx, fmt.Errorf("Trying to delete from cart more then cart have"))
	}

	err = m.DB.DeleteFromCartDB(ctx, user, sku, count)

	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil
}
