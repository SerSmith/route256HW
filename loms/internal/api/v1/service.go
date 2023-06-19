package service

import (
	"context"
	"route256/loms/internal/domain"
	"route256/loms/pkg/loms_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	loms_v1.UnimplementedLomsServer
	Model *domain.Model
}

func NewServer(Model *domain.Model) *service {
	return &service{Model: Model}

}

func (s *service) CancelOrder(ctx context.Context, in *loms_v1.CancelOrderRequest) (*emptypb.Empty, error) {

	OrderID := in.GetOrderID()
	err := s.Model.CancelOrder(ctx, OrderID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *service) CreateOrder(ctx context.Context, in *loms_v1.CreateOrderRequest) (*loms_v1.CreateOrderResponse, error) {

	Items := in.GetItems()
	User := in.GetUser()

	items := make([]domain.ItemOrder, 0, len(Items))
	for _, item := range Items {
		items = append(items, domain.ItemOrder{
			SKU:   item.SKU,
			Count: uint16(item.Count),
		})
	}
	id, err := s.Model.CreateOrder(ctx, User, items)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &loms_v1.CreateOrderResponse{
		OrderID: id,
	}, err
}

func (s *service) ListOrder(ctx context.Context, in *loms_v1.ListOrderRequest) (*loms_v1.ListOrderResponse, error) {

	OrderID := in.GetOrderID()

	order, err := s.Model.ListOrder(ctx, OrderID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	items := make([]*loms_v1.ItemOrder, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, &loms_v1.ItemOrder{
			SKU:   item.SKU,
			Count: uint32(item.Count),
		})
	}
	return &loms_v1.ListOrderResponse{
		User:   order.User,
		Items:  items,
		Status: string(order.Status),
	}, nil
}
func (s *service) OrderPayed(ctx context.Context, in *loms_v1.OrderPayedRequest) (*emptypb.Empty, error) {
	OrderID := in.GetOrderID()
	err := s.Model.OrderPayed(ctx, OrderID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *service) Stocks(ctx context.Context, in *loms_v1.StocksRequest) (*loms_v1.StocksResponse, error) {
	SKU := in.GetSKU()
	stocks, err := s.Model.Stocks(ctx, SKU)
	if err != nil {

		return nil, err
	}
	respStocks := make([]*loms_v1.ItemStock, 0, len(stocks))
	for _, stock := range stocks {
		respStocks = append(respStocks, &loms_v1.ItemStock{
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}
	return &loms_v1.StocksResponse{
		Stocks: respStocks,
	}, nil
}
