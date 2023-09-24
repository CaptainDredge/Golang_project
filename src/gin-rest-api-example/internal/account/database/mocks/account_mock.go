// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "gin-rest-api-example/internal/account/model"
)

// AccountDB is an autogenerated mock type for the AccountDB type
type AccountDB struct {
	mock.Mock
}

// FindByEmail provides a mock function with given fields: ctx, email
func (_m *AccountDB) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	ret := _m.Called(ctx, email)

	var r0 *model.Account
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Account); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, account
func (_m *AccountDB) Save(ctx context.Context, account *model.Account) error {
	ret := _m.Called(ctx, account)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Account) error); ok {
		r0 = rf(ctx, account)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, email, account
func (_m *AccountDB) Update(ctx context.Context, email string, account *model.Account) error {
	ret := _m.Called(ctx, email, account)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Account) error); ok {
		r0 = rf(ctx, email, account)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAccountDB interface {
	mock.TestingT
	Cleanup(func())
}

// NewAccountDB creates a new instance of AccountDB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAccountDB(t mockConstructorTestingTNewAccountDB) *AccountDB {
	mock := &AccountDB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
