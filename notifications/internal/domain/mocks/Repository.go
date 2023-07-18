// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	context "context"
	domain "route256/notifications/internal/domain"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// ReadNotifications provides a mock function with given fields: ctx, req
func (_m *Repository) ReadNotifications(ctx context.Context, req domain.NotificationHistoryRequest) ([]domain.NotificationMem, error) {
	ret := _m.Called(ctx, req)

	var r0 []domain.NotificationMem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.NotificationHistoryRequest) ([]domain.NotificationMem, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.NotificationHistoryRequest) []domain.NotificationMem); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.NotificationMem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.NotificationHistoryRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WriteNotification provides a mock function with given fields: ctx, message, DT
func (_m *Repository) WriteNotification(ctx context.Context, message domain.StatusChangeMessage, DT time.Time) error {
	ret := _m.Called(ctx, message, DT)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.StatusChangeMessage, time.Time) error); ok {
		r0 = rf(ctx, message, DT)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}