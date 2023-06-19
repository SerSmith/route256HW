package postgres

import (
	"context"
	"fmt"
	"log"
	"route256/checkout/internal/converter/schema2domain"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository/schema"
	"route256/libs/tx"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	provider tx.DBProvider
}

func New(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}

const (
	tableCart = "cart"
)

func (r *Repository) AddToCartDB(ctx context.Context, user int64, sku uint32, count uint16) error {

	db := r.provider.GetDB(ctx)
	query := `INSERT INTO cart (user_id, sku, "count") VALUES 
	($1, $2, $3)
ON CONFLICT (user_id, sku) DO UPDATE 
	SET "count"=cart.count + $3;`

	_, err := db.Exec(ctx, query, user, sku, count)

	if err != nil {
		return fmt.Errorf("exec insert stocks: %v", err)
	}

	return nil
}

func (r *Repository) DeleteFromCartDB(ctx context.Context, user int64, sku uint32, count uint16) error {

	db := r.provider.GetDB(ctx)

	query := `UPDATE cart
	SET "count"=cart.count - $3
	WHERE user_id = $1 and sku = $2;`

	_, err := db.Exec(ctx, query, user, sku, count)

	if err != nil {
		return fmt.Errorf("exec insert stocks: %v", err)
	}

	query = `DELETE
	FROM cart
	WHERE "count" <= 0`

	_, err = db.Exec(ctx, query)

	if err != nil {
		return fmt.Errorf("delete 0 rows: %v", err)
	}

	return nil
}

func (r *Repository) GetCartQauntDB(ctx context.Context, user int64, sku uint32) (uint16, error) {

	db := r.provider.GetDB(ctx)

	query := psql.Select("Count").
		From(tableCart).
		Where(sq.Eq{"user_id": user, "sku": sku})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query for ReservePtoduct get: %s", err)
	}

	var count []uint16
	err = pgxscan.Select(ctx, db, &count, rawSQL, args...)

	/* Если мы ничего не нашли, то у нас 0 соответствующих строчек*/
	if len(count) == 0 {
		count = append(count, 0)
	}

	if err != nil {
		return 0, fmt.Errorf("exec for ReservePtoduct get: %w", err)
	}

	return count[0], nil
}

func (r *Repository) GetCartDB(ctx context.Context, user int64) ([]domain.ItemOrder, error) {

	db := r.provider.GetDB(ctx)

	query := psql.Select("sku", "Count").
		From(tableCart).
		Where(sq.Eq{"user_id": user})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for ReservePtoduct get: %s", err)
	}

	var itemsInCartSchema []schema.ItemOrder
	err = pgxscan.Select(ctx, db, &itemsInCartSchema, rawSQL, args...)

	if err != nil {
		return nil, fmt.Errorf("exec for ReservePtoduct get: %w", err)
	}

	var itemsInCart []domain.ItemOrder

	for _, item := range itemsInCartSchema {
		itemsInCart = append(itemsInCart, schema2domain.ItemOrderConvert(item))
	}

	return itemsInCart, nil
}

func (r *Repository) WipeCartDB(ctx context.Context, user int64) error {
	db := r.provider.GetDB(ctx)

	query := psql.Delete(tableCart).
		Where(sq.Eq{"user_id": user})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for WipeCartDB get: %w", err)
	}

	log.Print(rawSQL, " ", args)

	_, err = db.Exec(ctx, rawSQL, args...)

	if err != nil {
		return fmt.Errorf("exec for ReservePtoduct get: %w", err)
	}

	return nil
}

func (r *Repository) RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error {
	return r.provider.RunRepeatableRead(ctx, fn)
}
