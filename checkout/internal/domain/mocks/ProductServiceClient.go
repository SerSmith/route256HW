// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "route256/checkout/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// ProductServiceClient is an autogenerated mock type for the ProductServiceClient type
type ProductServiceClient struct {
	mock.Mock
}

// GetProduct provides a mock function with given fields: ctx, sku
func (_m *ProductServiceClient) GetProduct(ctx context.Context, sku uint32) (*domain.Product, error) {
	ret := _m.Called(ctx, sku)

	var r0 *domain.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) (*domain.Product, error)); ok {
		return rf(ctx, sku)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint32) *domain.Product); ok {
		r0 = rf(ctx, sku)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint32) error); ok {
		r1 = rf(ctx, sku)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProductServiceClient creates a new instance of ProductServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductServiceClient {
	mock := &ProductServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
