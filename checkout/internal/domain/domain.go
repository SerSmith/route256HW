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

type Repository interface {
	AddToCartDB(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFromCartDB(ctx context.Context, user int64, sku uint32, count uint16) error
	GetCartQauntDB(ctx context.Context, user int64, sku uint32) (uint16, error)
	GetCartDB(ctx context.Context, user int64) ([]ItemOrder, error)
	WipeCartDB(ctx context.Context, user int64) error
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
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
	DB                   Repository
}

func New(LomsClient LomsClient, productServiceClient ProductServiceClient, DB Repository) *Model {
	return &Model{
		LomsClient:           LomsClient,
		productServiceClient: productServiceClient,
		DB:                   DB,
	}
}
