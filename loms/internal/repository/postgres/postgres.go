package postgres

import (
	"context"
	"fmt"
	"route256/libs/tx"
	"route256/loms/internal/domain"
	"github.com/georgysavva/scany/pgxscan"
	sq "github.com/Masterminds/squirrel"
	"log"

)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	provider tx.DBProvider
}

func New(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}

const (
	tableNameOrderUsers		= "ORDERUSERS"
	tableNameOrders			= "ORDERS"
	tableOrdersStatus	    = "OrdersStatus"
	tableStocksAvalible		= "StocksAvalible"
	tableStocksBought	  	= "StocksBought"
	tableStocksReserved		= "StocksReserved"
)


func (r *Repository) WriteOrderUser(ctx context.Context, User int64) (int64, error) {
	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableNameOrderUsers).Columns("user_id").Values(User).Suffix("RETURNING orderid")

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query for create orderuser: %s", err)
	}

	log.Println("WriteOrderUser", rawSQL, args)
	var orderID int64
	err = db.QueryRow(ctx, rawSQL, args...).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("exec insert orderuser: %s", err)
	}

	return orderID, err
}

func (r *Repository) WriteOrderItems(ctx context.Context, items []domain.ItemOrder, orderID int64) error {
	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableNameOrders).Columns("orderID", "SKU", "Count").Suffix("RETURNING orderid")

	for _, item := range items {
		query = query.
			Values(orderID, item.SKU, item.Count)

	}

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for create order: %s", err)
	}


	var out int 
	err = db.QueryRow(ctx, rawSQL, args...).Scan(&out)
	if err != nil {
		return fmt.Errorf("exec insert order: %w", err)
	}


	return nil
}


func (r *Repository) ChangeOrderStatus(ctx context.Context, orderID int64, Status domain.OrderStatus) (error) {
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
		return fmt.Errorf("exec insert order: %w", err)
	}



	return nil

}

func (r *Repository) ReserveProducts(ctx context.Context, orderID int64, stockInfos []domain.StockInfo) (error){

	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableStocksReserved).Columns("orderID", "sku", "warehouseID", "count")

	for _, stockInfo := range stockInfos {
		query = query.
			Values(orderID, stockInfo.SKU, stockInfo.WarehouseID, stockInfo.Count)

	}
	
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for ReserveProducts: %s", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec ReserveProducts: %w", err)
	}

	return nil
}



func (r *Repository) UnreserveProducts(ctx context.Context, orderID int64) (error){

	db := r.provider.GetDB(ctx)

	query := psql.Delete(tableStocksReserved).
			 Where(sq.Eq{"orderID": orderID})


	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for ReserveProducts: %s", err)
	}

	_, err = db.Query(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec ReserveProducts: %w", err)
	}

	return nil
}

func (r *Repository) BuyProducts(ctx context.Context, stocks []domain.StockInfo) (error){

	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableStocksBought).Columns( "sku", "warehouseID", "count")

	for _, stock := range stocks {
		query = query.
			Values(stock.SKU, stock.WarehouseID, stock.Count)

	}
	
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for ReserveProducts: %s", err)
	}

	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec ReserveProducts: %w", err)
	}

	return nil
}




func (r *Repository) MinusAvalibleCount(ctx context.Context, stockInfos[]domain.StockInfo) (error){

	db := r.provider.GetDB(ctx)


	query := `INSERT INTO StocksAvalible (sku, warehouseID, "count") VALUES 
	($1, $2, $3)
ON CONFLICT (sku, warehouseID) DO UPDATE 
	SET "count"=StocksAvalible.count - $3;`


	for _, stockInfo := range stockInfos{

		tmp, err := db.Exec(ctx, query, stockInfo.SKU, stockInfo.WarehouseID, stockInfo.Count)

		log.Println("tmp ", tmp)

		if err != nil {
			return fmt.Errorf("exec insert stocks: %v", err)
		}

	}


	return nil
}

func (r *Repository) PlusAvalibleCount(ctx context.Context, stockInfos[]domain.StockInfo) (error){

	db := r.provider.GetDB(ctx)
	query := `INSERT INTO StocksAvalible (sku, warehouseID, "count") VALUES 
	($1, $2, $3)
ON CONFLICT (sku, warehouseID) DO UPDATE 
	SET "count"=StocksAvalible.count + $3;`


	for _, stockInfo := range stockInfos{


		_, err := db.Exec(ctx, query, stockInfo.SKU, stockInfo.WarehouseID, stockInfo.Count)

		if err != nil {
			return fmt.Errorf("exec insert stocks: %v", err)
		}



	}

	return nil
}

func (r *Repository) GetAvailableBySku(ctx context.Context, sku uint32) ([]domain.Stock, error){
	db := r.provider.GetDB(ctx)

	query := psql.Select("WarehouseID", "Count").
	From(tableStocksAvalible).
	Where(sq.Eq{"sku": sku})


	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for ReservePtoduct get: %s", err)
	}


	var stocksFound []domain.Stock
	err = pgxscan.Select(ctx, db, &stocksFound, rawSQL, args...)

	if err != nil {
		return nil, fmt.Errorf("exec for ReservePtoduct get: %w", err)
	}

	return stocksFound, nil
}



func (r *Repository) GetReservedByOrderID(ctx context.Context, orderID int64) ([]domain.StockInfo, error){
	db := r.provider.GetDB(ctx)

	query := psql.Select("sku", "WarehouseID", "Count").
	From(tableStocksReserved).
	Where(sq.Eq{"orderID": orderID})




	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for ReservePtoduct get: %s", err)
	}


	var stocks []domain.StockInfo
	err = pgxscan.Select(ctx, db, &stocks, rawSQL, args...)

	if err != nil {
		return nil, fmt.Errorf("exec for ReservePtoduct get: %w", err)
	}

	return stocks, nil
}






func (r *Repository) GetOrderDetails(ctx context.Context, orderID int64) (domain.Order, error){
	
	db := r.provider.GetDB(ctx)



	query := psql.Select("user_id").
	From(tableNameOrderUsers).
	Where(sq.Eq{"orderID": orderID})



	rawSQL, args, err := query.ToSql()
	if err != nil {
		return domain.Order{}, fmt.Errorf("build query for ReservePtoduct get: %s", err)
	}


	var userID []int64
	err = pgxscan.Select(ctx, db, &userID, rawSQL, args...)


	if err != nil {
		return domain.Order{}, fmt.Errorf("exec for ReservePtoduct get: %w", err)
	}




	query = psql.Select("status").
	From(tableOrdersStatus).
	Where(sq.Eq{"orderID": orderID})



	rawSQL, args, err = query.ToSql()
	if err != nil {
		return domain.Order{}, fmt.Errorf("build query for ReservePtoduct get: %s", err)
	}


	var status []domain.OrderStatus
	err = pgxscan.Select(ctx, db, &status, rawSQL, args...)

	if err != nil {
		return domain.Order{}, fmt.Errorf("exec for ReservePtoduct get: %w", err)
	}





	query = psql.Select("sku", "count").
	From(tableNameOrders).
	Where(sq.Eq{"orderID": orderID})



	rawSQL, args, err = query.ToSql()
	if err != nil {
		return domain.Order{}, fmt.Errorf("build query for ReservePtoduct get: %s", err)
	}


	var items []domain.ItemOrder
	err = pgxscan.Select(ctx, db, &items, rawSQL, args...)

	if err != nil {
		return domain.Order{}, fmt.Errorf("exec for ReservePtoduct get: %w", err)
	}



	var itempointers []*domain.ItemOrder

	for _, i := range items{

		itempointers = append(itempointers, &i)

	}


	out := domain.Order{
			User: userID[0],
			Items: itempointers,
			Status: status[0],
		}

	return out, nil
}


func (r *Repository) GetOrderStatus(ctx context.Context, orderID int64) (domain.OrderStatus, error){
	db := r.provider.GetDB(ctx)

	query := psql.Select("Status").
	From(tableOrdersStatus).
	Where(sq.Eq{"orderID": orderID})


	rawSQL, args, err := query.ToSql()
	if err != nil {
		return "", fmt.Errorf("build query for ReservePtoduct get: %s", err)
	}


	var  status []domain.OrderStatus
	err = pgxscan.Select(ctx, db, &status, rawSQL, args...)

	if err != nil {
		return "", fmt.Errorf("exec for ReservePtoduct get: %w", err)
	}

	return status[0], nil
}