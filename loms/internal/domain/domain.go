package domain

import "context"



type Repository interface {
	WriteOrderUser(ctx context.Context, User int64) (int64, error)
	WriteOrderItems(ctx context.Context, items []ItemOrder, orderID int64) error
	ReserveProducts(ctx context.Context, orderID int64, stockInfos []StockInfo) (error)
	MinusAvalibleCount(ctx context.Context, stockInfos []StockInfo) (error)
	PlusAvalibleCount(ctx context.Context, stockInfos []StockInfo) (error)
	GetAvailableBySku(ctx context.Context, sku uint32) ([]Stock, error)
	ChangeOrderStatus(ctx context.Context, orderID int64, Status OrderStatus) (error)
	GetOrderStatus(ctx context.Context, orderID int64) (OrderStatus, error)
	GetOrderDetails(ctx context.Context, orderID int64) (Order, error)
	UnreserveProducts(ctx context.Context, orderID int64) (error)
	BuyProducts(ctx context.Context, stocks []StockInfo) (error)
	GetReservedByOrderID(ctx context.Context, orderID int64) ([]StockInfo, error)

}

type OrderStatus string

const (
	NewStatus             OrderStatus = "new"
	AwaitingPaymentStatus OrderStatus = "awaiting payment"
	FailedStatus          OrderStatus = "failed"
	PayedStatus           OrderStatus = "payed"
	CanceledStatus        OrderStatus = "cancelled"
)

type ItemOrder struct {
	SKU   uint32 `db:"sku"`
	Count uint16 `db:"count"`
}


type Order struct {
	User   int64
	Items  []*ItemOrder
	Status OrderStatus
}


type Stock struct {
	WarehouseID int64 `db:"warehouseid"`
	Count       uint64 `db:"count"`
}

type StockInfo struct {
	SKU			int64 `db:"sku"`
	WarehouseID int64 `db:"warehouseid"`
	Count       uint64 `db:"count"`
}

type Model struct {
	DB Repository
}

func New(DB Repository) *Model {
	return &Model{DB: DB}
}

