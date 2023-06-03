package service

import (

	"context"
	"route256/checkout/internal/domain"
	"route256/checkout/pkg/checkout_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)




type service struct {
	checkout_v1.UnimplementedCheckoutServer
	Model *domain.Model
}



func NewServer(loms_client domain.LomsClient, product_service_client domain.ProductServiceClient) *service {
	return &service{
					Model: domain.New(loms_client, product_service_client),
				}

}

func (s *service) AddToCart(ctx context.Context, in *checkout_v1.AddToCartRequest) (*emptypb.Empty, error) {

	User := in.GetUser()
	SKU := in.GetSku()
	Count := uint16(in.GetCount())


	err := s.Model.AddToCart(ctx, User, SKU, Count)
	
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}


func (s *service) DeleteFromCart(ctx context.Context, in *checkout_v1.DeleteFromCartRequest) (*emptypb.Empty, error) {

	User := in.GetUser()
	SKU := in.GetSku()
	Count := uint16(in.GetCount())

	err := s.Model.DeleteFromCart(ctx, User, SKU, Count)

	
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *service) ListCart(ctx context.Context, in *checkout_v1.ListCartRequest) (*checkout_v1.ListCartResponse, error) {

	User := in.GetUser()
	totalPrice, cart, err := s.Model.ListCart(ctx, User)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	items := make([]*checkout_v1.ItemCart, 0, len(cart))
	for _, item := range cart {
		items = append(items, &checkout_v1.ItemCart{
			SKU:   item.SKU,
			Count: uint32(item.Count),
			Price: item.Product.Price,
			Name:  item.Product.Name,
		})
	}

	resp := &checkout_v1.ListCartResponse{
		TotalPrice: totalPrice,
		Items:      items,
	}

	return resp, nil
}

func (s *service) Purchase(ctx context.Context, in *checkout_v1.PurchaseRequest) (*checkout_v1.PurchaseResponse, error) {

	User := in.GetUser()

	orderID, err := s.Model.Purchase(ctx, User)
	if err != nil {
		return nil, err
	}

	return &checkout_v1.PurchaseResponse{OrderID: orderID}, nil
}
