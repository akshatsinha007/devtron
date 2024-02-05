// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	bean "github.com/devtron-labs/devtron/pkg/globalPolicy/bean"
	mock "github.com/stretchr/testify/mock"

	pg "github.com/go-pg/pg"

	repository "github.com/devtron-labs/devtron/pkg/globalPolicy/repository"
)

// GlobalPolicyRepository is an autogenerated mock type for the GlobalPolicyRepository type
type GlobalPolicyRepository struct {
	mock.Mock
}

// CommitTransaction provides a mock function with given fields: tx
func (_m *GlobalPolicyRepository) CommitTransaction(tx *pg.Tx) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pg.Tx) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: model
func (_m *GlobalPolicyRepository) Create(model *repository.GlobalPolicy) error {
	ret := _m.Called(model)

	var r0 error
	if rf, ok := ret.Get(0).(func(*repository.GlobalPolicy) error); ok {
		r0 = rf(model)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllByPolicyOfAndVersion provides a mock function with given fields: policyOf, policyVersion
func (_m *GlobalPolicyRepository) GetAllByPolicyOfAndVersion(policyOf bean.GlobalPolicyType, policyVersion bean.GlobalPolicyVersion) ([]*repository.GlobalPolicy, error) {
	ret := _m.Called(policyOf, policyVersion)

	var r0 []*repository.GlobalPolicy
	if rf, ok := ret.Get(0).(func(bean.GlobalPolicyType, bean.GlobalPolicyVersion) []*repository.GlobalPolicy); ok {
		r0 = rf(policyOf, policyVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.GlobalPolicy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bean.GlobalPolicyType, bean.GlobalPolicyVersion) error); ok {
		r1 = rf(policyOf, policyVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id
func (_m *GlobalPolicyRepository) GetById(id int) (*repository.GlobalPolicy, error) {
	ret := _m.Called(id)

	var r0 *repository.GlobalPolicy
	if rf, ok := ret.Get(0).(func(int) *repository.GlobalPolicy); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.GlobalPolicy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByName provides a mock function with given fields: name
func (_m *GlobalPolicyRepository) GetByName(name string) (*repository.GlobalPolicy, error) {
	ret := _m.Called(name)

	var r0 *repository.GlobalPolicy
	if rf, ok := ret.Get(0).(func(string) *repository.GlobalPolicy); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.GlobalPolicy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDbTransaction provides a mock function with given fields:
func (_m *GlobalPolicyRepository) GetDbTransaction() (*pg.Tx, error) {
	ret := _m.Called()

	var r0 *pg.Tx
	if rf, ok := ret.Get(0).(func() *pg.Tx); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pg.Tx)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetEnabledPoliciesByIds provides a mock function with given fields: ids
func (_m *GlobalPolicyRepository) GetEnabledPoliciesByIds(ids []int) ([]*repository.GlobalPolicy, error) {
	ret := _m.Called(ids)

	var r0 []*repository.GlobalPolicy
	if rf, ok := ret.Get(0).(func([]int) []*repository.GlobalPolicy); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.GlobalPolicy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]int) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarkDeletedById provides a mock function with given fields: id, userId, tx
func (_m *GlobalPolicyRepository) MarkDeletedById(id int, userId int32, tx *pg.Tx) error {
	ret := _m.Called(id, userId, tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, int32, *pg.Tx) error); ok {
		r0 = rf(id, userId, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RollBackTransaction provides a mock function with given fields: tx
func (_m *GlobalPolicyRepository) RollBackTransaction(tx *pg.Tx) error {
	ret := _m.Called(tx)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pg.Tx) error); ok {
		r0 = rf(tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: model
func (_m *GlobalPolicyRepository) Update(model *repository.GlobalPolicy) error {
	ret := _m.Called(model)

	var r0 error
	if rf, ok := ret.Get(0).(func(*repository.GlobalPolicy) error); ok {
		r0 = rf(model)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewGlobalPolicyRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewGlobalPolicyRepository creates a new instance of GlobalPolicyRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewGlobalPolicyRepository(t mockConstructorTestingTNewGlobalPolicyRepository) *GlobalPolicyRepository {
	mock := &GlobalPolicyRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}