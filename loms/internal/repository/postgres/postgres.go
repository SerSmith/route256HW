package postgres

import (
	"context"
	"fmt"
	"route256/libs/tracer"
	"route256/libs/tx"
	"route256/loms/internal/converter/schema2domain"
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	provider tx.DBProvider
}

func New(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}

const (
	tableNameOrders     = "ORDERS"
	tableOrdersStatus   = "OrdersStatus"
	tableStocksAvalible = "StocksAvalible"
	tableStocksBought   = "StocksBought"
	tableStocksReserved = "StocksReserved"
)

func (r *Repository) WriteOrder(ctx context.Context, items []domain.ItemOrder, User int64) (int64, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/WriteOrder")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableNameOrders).Columns("user_id", "SKU", "Count").Suffix("RETURNING orderid")

	for _, item := range items {
		query = query.
			Values(User, item.SKU, item.Count)

	}

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, err)
	}

	var orderID int64
	err = db.QueryRow(ctx, rawSQL, args...).Scan(&orderID)
	if err != nil {
		return 0, tracer.MarkSpanWithError(ctx, err)
	}

	return orderID, nil
}

func (r *Repository) ChangeOrderStatus(ctx context.Context, orderID int64, Status domain.OrderStatus) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/ChangeOrderStatus")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := fmt.Sprintf(`INSERT INTO OrdersStatus (orderID, Status) VALUES 
							(%d, '%s')
							ON CONFLICT (orderID) DO UPDATE 
							SET Status='%s'
							RETURNING Status;
							`,
		orderID, Status, Status)

	var out string
	err := db.QueryRow(ctx, query).Scan(&out)

	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil

}

func (r *Repository) ReserveProducts(ctx context.Context, orderID int64, stockInfos []domain.StockInfo) error {

	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/ReserveProducts")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableStocksReserved).Columns("orderID", "sku", "warehouseID", "count")

	for _, stockInfo := range stockInfos {
		query = query.
			Values(orderID, stockInfo.SKU, stockInfo.WarehouseID, stockInfo.Count)

	}

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil
}

func (r *Repository) UnreserveProducts(ctx context.Context, orderID int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/UnreserveProducts")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := psql.Delete(tableStocksReserved).
		Where(sq.Eq{"orderID": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	_, err = db.Query(ctx, rawSQL, args...)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil
}

func (r *Repository) BuyProducts(ctx context.Context, stocks []domain.StockInfo) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/BuyProducts")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableStocksBought).Columns("sku", "warehouseID", "count")

	for _, stock := range stocks {
		query = query.
			Values(stock.SKU, stock.WarehouseID, stock.Count)

	}

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	return nil
}

func (r *Repository) MinusAvalibleCount(ctx context.Context, stockInfos []domain.StockInfo) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/MinusAvalibleCount")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := `INSERT INTO StocksAvalible (sku, warehouseID, "count") VALUES 
	($1, $2, $3)
ON CONFLICT (sku, warehouseID) DO UPDATE 
	SET "count"=StocksAvalible.count - $3;`

	for _, stockInfo := range stockInfos {

		_, err := db.Exec(ctx, query, stockInfo.SKU, stockInfo.WarehouseID, stockInfo.Count)

		if err != nil {
			return tracer.MarkSpanWithError(ctx, err)
		}

	}

	return nil
}

func (r *Repository) PlusAvalibleCount(ctx context.Context, stockInfos []domain.StockInfo) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/PlusAvalibleCount")
	defer span.Finish()

	db := r.provider.GetDB(ctx)
	query := `INSERT INTO StocksAvalible (sku, warehouseID, "count") VALUES 
	($1, $2, $3)
ON CONFLICT (sku, warehouseID) DO UPDATE 
	SET "count"=StocksAvalible.count + $3;`

	for _, stockInfo := range stockInfos {

		_, err := db.Exec(ctx, query, stockInfo.SKU, stockInfo.WarehouseID, stockInfo.Count)

		if err != nil {
			return tracer.MarkSpanWithError(ctx, err)
		}

	}

	return nil
}

func (r *Repository) GetAvailableBySku(ctx context.Context, sku uint32) ([]domain.Stock, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/GetAvailableBySku")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := psql.Select("WarehouseID", "Count").
		From(tableStocksAvalible).
		Where(sq.Eq{"sku": sku})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, err)
	}

	var stocksFoundSchema []schema.Stock
	err = pgxscan.Select(ctx, db, &stocksFoundSchema, rawSQL, args...)

	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, err)
	}

	var stocksFound []domain.Stock

	for _, stock := range stocksFoundSchema {
		stocksFound = append(stocksFound, schema2domain.StockConvert(stock))
	}

	return stocksFound, nil
}

func (r *Repository) GetReservedByOrderID(ctx context.Context, orderID int64) ([]domain.StockInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/GetReservedByOrderID")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := psql.Select("sku", "WarehouseID", "Count").
		From(tableStocksReserved).
		Where(sq.Eq{"orderID": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, err)
	}

	var stocksSchema []schema.StockInfo
	err = pgxscan.Select(ctx, db, &stocksSchema, rawSQL, args...)

	if err != nil {
		return nil, tracer.MarkSpanWithError(ctx, err)
	}

	var stocks []domain.StockInfo

	for _, stock := range stocksSchema {
		stocks = append(stocks, schema2domain.StockInfoConvert(stock))
	}

	return stocks, nil
}

func (r *Repository) GetOrderDetails(ctx context.Context, orderID int64) (domain.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/GetOrderDetails")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := psql.Select("status").
		From(tableOrdersStatus).
		Where(sq.Eq{"orderID": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return domain.Order{}, tracer.MarkSpanWithError(ctx, err)
	}

	var status []domain.OrderStatus
	err = pgxscan.Select(ctx, db, &status, rawSQL, args...)

	if err != nil {
		return domain.Order{}, tracer.MarkSpanWithError(ctx, err)
	}

	query = psql.Select("user_id").
		From(tableNameOrders).
		Where(sq.Eq{"orderID": orderID})

	rawSQL, args, err = query.ToSql()
	if err != nil {
		return domain.Order{}, tracer.MarkSpanWithError(ctx, err)
	}

	var userID []int64
	err = pgxscan.Select(ctx, db, &userID, rawSQL, args...)

	if err != nil {
		return domain.Order{}, tracer.MarkSpanWithError(ctx, err)
	}

	query = psql.Select("sku", "count").
		From(tableNameOrders).
		Where(sq.Eq{"orderID": orderID})

	rawSQL, args, err = query.ToSql()
	if err != nil {
		return domain.Order{}, tracer.MarkSpanWithError(ctx, err)
	}

	var itemsSchema []schema.ItemOrder
	err = pgxscan.Select(ctx, db, &itemsSchema, rawSQL, args...)

	if err != nil {
		return domain.Order{}, tracer.MarkSpanWithError(ctx, err)
	}

	var items []domain.ItemOrder
	for _, item := range itemsSchema {
		items = append(items, schema2domain.ItemOrderConvert(item))
	}

	var itempointers []*domain.ItemOrder

	for _, i := range items {

		itempointers = append(itempointers, &i)

	}

	out := domain.Order{
		User:   userID[0],
		Items:  itempointers,
		Status: status[0],
	}

	return out, nil
}

func (r *Repository) GetOrderStatus(ctx context.Context, orderID int64) (domain.OrderStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repository/postgres/GetOrderStatus")
	defer span.Finish()

	db := r.provider.GetDB(ctx)

	query := psql.Select("Status").
		From(tableOrdersStatus).
		Where(sq.Eq{"orderID": orderID})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return "", tracer.MarkSpanWithError(ctx, err)
	}

	var status []domain.OrderStatus
	err = pgxscan.Select(ctx, db, &status, rawSQL, args...)

	if err != nil {
		return "", tracer.MarkSpanWithError(ctx, err)
	}

	var out domain.OrderStatus

	if len(status) > 0 {
		out = status[0]
	} else {
		out = domain.NullStatus
	}

	return out, nil
}

func (r *Repository) RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error {
	return r.provider.RunRepeatableRead(ctx, fn)
}
