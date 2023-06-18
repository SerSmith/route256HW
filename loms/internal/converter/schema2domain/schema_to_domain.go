package schema2domain

import (
	"route256/loms/internal/domain"
	"route256/loms/internal/repository/schema"
)

func ItemOrderConvert(in schema.ItemOrder) domain.ItemOrder {
	return domain.ItemOrder{
		SKU:   in.SKU,
		Count: in.Count}
}

func StockConvert(in schema.Stock) domain.Stock {
	return domain.Stock{
		WarehouseID: in.WarehouseID,
		Count:       in.Count}
}

func StockInfoConvert(in schema.StockInfo) domain.StockInfo {
	return domain.StockInfo{
		SKU:         in.SKU,
		Count:       in.Count,
		WarehouseID: in.WarehouseID}
}
