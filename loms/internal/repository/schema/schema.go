package schema

type ItemOrder struct {
	SKU   uint32 `db:"sku"`
	Count uint16 `db:"count"`
}

type Stock struct {
	WarehouseID int64  `db:"warehouseid"`
	Count       uint64 `db:"count"`
}

type StockInfo struct {
	SKU         int64  `db:"sku"`
	WarehouseID int64  `db:"warehouseid"`
	Count       uint64 `db:"count"`
}
