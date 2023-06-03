package domain

import (
	"context"
)

type LomsClient interface {
	Stocks(ctx context.Context, sku uint32) ([]Stock, error)
	CreateOrder(ctx context.Context, user int64, items []ItemOrder) (int64, error)
}

type ProductServiceClient interface {
	GetProduct(ctx context.Context, sku uint32) (*Product, error)
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type Product struct {
	Name  string
	Price uint32
}

type ItemCart struct {
	SKU     uint32
	Count   uint16
	Product Product
}

type ItemOrder struct {
	SKU   uint32
	Count uint16
}

type Model struct {
	LomsClient           LomsClient
	productServiceClient ProductServiceClient
}

func New(LomsClient LomsClient, productServiceClient ProductServiceClient) *Model {
	return &Model{
		LomsClient:           LomsClient,
		productServiceClient: productServiceClient,
	}
}
