// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/imbonda/vmm-bot/pkg/models"
)

// ExchangeClient is an autogenerated mock type for the ExchangeClient type
type ExchangeClient struct {
	mock.Mock
}

// GetOrderBook provides a mock function with given fields: ctx, symbol
func (_m *ExchangeClient) GetOrderBook(ctx context.Context, symbol string) (*models.OrderBook, error) {
	ret := _m.Called(ctx, symbol)

	if len(ret) == 0 {
		panic("no return value specified for GetOrderBook")
	}

	var r0 *models.OrderBook
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.OrderBook, error)); ok {
		return rf(ctx, symbol)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.OrderBook); ok {
		r0 = rf(ctx, symbol)
	} else if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.OrderBook)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, symbol)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlaceOrder provides a mock function with given fields: ctx, order
func (_m *ExchangeClient) PlaceOrder(ctx context.Context, order *models.Order) error {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for PlaceOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewExchangeClient creates a new instance of ExchangeClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExchangeClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExchangeClient {
	mock := &ExchangeClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
