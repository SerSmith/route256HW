package schema2domain

import (
	"route256/checkout/internal/domain"
	"route256/checkout/internal/repository/schema"
)

func ItemOrderConvert(in schema.ItemOrder) domain.ItemOrder {
	return domain.ItemOrder{
		SKU:   in.SKU,
		Count: in.Count}
}
