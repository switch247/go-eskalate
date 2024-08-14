// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"
	Domain "main/Domain"

	mock "github.com/stretchr/testify/mock"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// CreateUsers provides a mock function with given fields: ctx, user
func (_m *UserRepository) CreateUsers(ctx context.Context, user *Domain.User) (Domain.OmitedUser, error, int) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUsers")
	}

	var r0 Domain.OmitedUser
	var r1 error
	var r2 int
	if rf, ok := ret.Get(0).(func(context.Context, *Domain.User) (Domain.OmitedUser, error, int)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *Domain.User) Domain.OmitedUser); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(Domain.OmitedUser)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *Domain.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	if rf, ok := ret.Get(2).(func(context.Context, *Domain.User) int); ok {
		r2 = rf(ctx, user)
	} else {
		r2 = ret.Get(2).(int)
	}

	return r0, r1, r2
}

// DeleteUsersById provides a mock function with given fields: ctx, id, user
func (_m *UserRepository) DeleteUsersById(ctx context.Context, id primitive.ObjectID, user Domain.OmitedUser) (error, int) {
	ret := _m.Called(ctx, id, user)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUsersById")
	}

	var r0 error
	var r1 int
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, Domain.OmitedUser) (error, int)); ok {
		return rf(ctx, id, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, Domain.OmitedUser) error); ok {
		r0 = rf(ctx, id, user)
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func(context.Context, primitive.ObjectID, Domain.OmitedUser) int); ok {
		r1 = rf(ctx, id, user)
	} else {
		r1 = ret.Get(1).(int)
	}

	return r0, r1
}

// GetUsers provides a mock function with given fields: ctx
func (_m *UserRepository) GetUsers(ctx context.Context) ([]*Domain.OmitedUser, error, int) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetUsers")
	}

	var r0 []*Domain.OmitedUser
	var r1 error
	var r2 int
	if rf, ok := ret.Get(0).(func(context.Context) ([]*Domain.OmitedUser, error, int)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*Domain.OmitedUser); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Domain.OmitedUser)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	if rf, ok := ret.Get(2).(func(context.Context) int); ok {
		r2 = rf(ctx)
	} else {
		r2 = ret.Get(2).(int)
	}

	return r0, r1, r2
}

// GetUsersById provides a mock function with given fields: ctx, id, user
func (_m *UserRepository) GetUsersById(ctx context.Context, id primitive.ObjectID, user Domain.OmitedUser) (Domain.OmitedUser, error, int) {
	ret := _m.Called(ctx, id, user)

	if len(ret) == 0 {
		panic("no return value specified for GetUsersById")
	}

	var r0 Domain.OmitedUser
	var r1 error
	var r2 int
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, Domain.OmitedUser) (Domain.OmitedUser, error, int)); ok {
		return rf(ctx, id, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, Domain.OmitedUser) Domain.OmitedUser); ok {
		r0 = rf(ctx, id, user)
	} else {
		r0 = ret.Get(0).(Domain.OmitedUser)
	}

	if rf, ok := ret.Get(1).(func(context.Context, primitive.ObjectID, Domain.OmitedUser) error); ok {
		r1 = rf(ctx, id, user)
	} else {
		r1 = ret.Error(1)
	}

	if rf, ok := ret.Get(2).(func(context.Context, primitive.ObjectID, Domain.OmitedUser) int); ok {
		r2 = rf(ctx, id, user)
	} else {
		r2 = ret.Get(2).(int)
	}

	return r0, r1, r2
}

// UpdateUsersById provides a mock function with given fields: ctx, id, user, curentuser
func (_m *UserRepository) UpdateUsersById(ctx context.Context, id primitive.ObjectID, user Domain.User, curentuser Domain.OmitedUser) (Domain.OmitedUser, error, int) {
	ret := _m.Called(ctx, id, user, curentuser)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUsersById")
	}

	var r0 Domain.OmitedUser
	var r1 error
	var r2 int
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, Domain.User, Domain.OmitedUser) (Domain.OmitedUser, error, int)); ok {
		return rf(ctx, id, user, curentuser)
	}
	if rf, ok := ret.Get(0).(func(context.Context, primitive.ObjectID, Domain.User, Domain.OmitedUser) Domain.OmitedUser); ok {
		r0 = rf(ctx, id, user, curentuser)
	} else {
		r0 = ret.Get(0).(Domain.OmitedUser)
	}

	if rf, ok := ret.Get(1).(func(context.Context, primitive.ObjectID, Domain.User, Domain.OmitedUser) error); ok {
		r1 = rf(ctx, id, user, curentuser)
	} else {
		r1 = ret.Error(1)
	}

	if rf, ok := ret.Get(2).(func(context.Context, primitive.ObjectID, Domain.User, Domain.OmitedUser) int); ok {
		r2 = rf(ctx, id, user, curentuser)
	} else {
		r2 = ret.Get(2).(int)
	}

	return r0, r1, r2
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
