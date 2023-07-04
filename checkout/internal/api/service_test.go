package service

import (
	"context"
	// "errors"
	"fmt"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/domain/mocks"
	"route256/checkout/pkg/checkout_v1"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_ListCart(t *testing.T) {
	t.Parallel()

	t.Run("success, 2 known items", func(t *testing.T) {
		t.Parallel()

		const (
			userID = int64(1)

			knownItem1 = uint32(1)
			knownItem2 = uint32(2)
		)

		var (
			user1 = checkout_v1.ListCartRequest{
				User: userID}

			Item1 = domain.ItemOrder{SKU: knownItem1,
				Count: 10}

			Item2 = domain.ItemOrder{SKU: knownItem2,
				Count: 20}

			Prodct1 = domain.Product{
				Name:  "1",
				Price: 1}

			Prodct2 = domain.Product{
				Name:  "2",
				Price: 2}

			ItemOut1 = checkout_v1.ItemCart{SKU: Item1.SKU,
				Count: uint32(Item1.Count),
				Name:  Prodct1.Name,
				Price: Prodct1.Price}
			ItemOut2 = checkout_v1.ItemCart{SKU: Item2.SKU,
				Count: uint32(Item2.Count),
				Name:  Prodct2.Name,
				Price: Prodct2.Price}

			correctAnswer = &checkout_v1.ListCartResponse{
				Items:      []*checkout_v1.ItemCart{&ItemOut1, &ItemOut2},
				TotalPrice: ItemOut1.Count*ItemOut1.Price + ItemOut2.Count*ItemOut2.Price}
		)

		repositoryMock := mocks.NewRepository(t)
		repositoryMock.On("GetCartDB", mock.Anything, userID).Return([]domain.ItemOrder{Item1, Item2}, nil).Once()

		productServiceClientMock := mocks.NewProductServiceClient(t)
		productServiceClientMock.On("GetProduct", mock.Anything, knownItem1).Return(&Prodct1, nil)
		productServiceClientMock.On("GetProduct", mock.Anything, knownItem2).Return(&Prodct2, nil)

		serv := NewServer(nil, productServiceClientMock, repositoryMock)

		// Act

		ans, err := serv.ListCart(context.Background(), &user1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, ans.TotalPrice, correctAnswer.TotalPrice)
		require.ElementsMatch(t, ans.Items, correctAnswer.Items)
	})

	t.Run("success, 2 known 1 unknown item", func(t *testing.T) {
		t.Parallel()

		const (
			userID = int64(1)

			knownItem1  = uint32(1)
			knownItem2  = uint32(2)
			unknownItem = uint32(3)
		)

		var (
			user1 = checkout_v1.ListCartRequest{
				User: userID}

			Item1 = domain.ItemOrder{SKU: knownItem1,
				Count: 10}

			Item2 = domain.ItemOrder{SKU: knownItem2,
				Count: 20}

			Item3 = domain.ItemOrder{SKU: unknownItem,
				Count: 30}

			Prodct1 = domain.Product{
				Name:  "1",
				Price: 1}

			Prodct2 = domain.Product{
				Name:  "2",
				Price: 2}

			ItemOut1 = checkout_v1.ItemCart{SKU: Item1.SKU,
				Count: uint32(Item1.Count),
				Name:  Prodct1.Name,
				Price: Prodct1.Price}
			ItemOut2 = checkout_v1.ItemCart{SKU: Item2.SKU,
				Count: uint32(Item2.Count),
				Name:  Prodct2.Name,
				Price: Prodct2.Price}
			ItemOut3 = checkout_v1.ItemCart{SKU: Item3.SKU,
				Count: uint32(Item3.Count),
				Name:  domain.GetProductUnknownName,
				Price: domain.GetProductUnknownPrice}

			correctAnswer = &checkout_v1.ListCartResponse{
				Items:      []*checkout_v1.ItemCart{&ItemOut1, &ItemOut2, &ItemOut3},
				TotalPrice: ItemOut1.Count*ItemOut1.Price + ItemOut2.Count*ItemOut2.Price}
		)

		repositoryMock := mocks.NewRepository(t)
		repositoryMock.On("GetCartDB", mock.Anything, userID).Return([]domain.ItemOrder{Item1, Item2, Item3}, nil).Once()

		productServiceClientMock := mocks.NewProductServiceClient(t)
		productServiceClientMock.On("GetProduct", mock.Anything, knownItem1).Return(&Prodct1, nil)
		productServiceClientMock.On("GetProduct", mock.Anything, knownItem2).Return(&Prodct2, nil)
		productServiceClientMock.On("GetProduct", mock.Anything, unknownItem).Return(nil, fmt.Errorf("error"))

		serv := NewServer(nil, productServiceClientMock, repositoryMock)

		// Act

		ans, err := serv.ListCart(context.Background(), &user1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, ans.TotalPrice, correctAnswer.TotalPrice)
		require.ElementsMatch(t, ans.Items, correctAnswer.Items)
	})

	t.Run("success, no items", func(t *testing.T) {
		t.Parallel()

		const (
			userID = int64(1)
		)

		var (
			user1 = checkout_v1.ListCartRequest{
				User: userID}

			correctAnswer = &checkout_v1.ListCartResponse{
				Items:      []*checkout_v1.ItemCart{},
				TotalPrice: 0}
		)

		repositoryMock := mocks.NewRepository(t)
		repositoryMock.On("GetCartDB", mock.Anything, userID).Return([]domain.ItemOrder{}, nil).Once()

		productServiceClientMock := mocks.NewProductServiceClient(t)

		serv := NewServer(nil, productServiceClientMock, repositoryMock)

		// Act

		ans, err := serv.ListCart(context.Background(), &user1)

		// Assert
		require.NoError(t, err)
		require.Equal(t, ans.TotalPrice, correctAnswer.TotalPrice)
		require.ElementsMatch(t, ans.Items, correctAnswer.Items)
	})

}
