// Code generated by mockery v2.36.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/RomanUtolin/RESTful-CRUD/internall/entity"
	mock "github.com/stretchr/testify/mock"
)

// PersonLogic is an autogenerated mock type for the PersonLogic type
type PersonLogic struct {
	mock.Mock
}

// Create provides a mock function with given fields: req
func (_m *PersonLogic) Create(req *entity.Person) (*entity.Person, error) {
	ret := _m.Called(req)

	var r0 *entity.Person
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.Person) (*entity.Person, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*entity.Person) *entity.Person); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Person)
		}
	}

	if rf, ok := ret.Get(1).(func(*entity.Person) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *PersonLogic) Delete(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: email, phone, firstName, limit, offset
func (_m *PersonLogic) GetAll(email string, phone string, firstName string, limit int, offset int) ([]*entity.Person, int, error) {
	ret := _m.Called(email, phone, firstName, limit, offset)

	var r0 []*entity.Person
	var r1 int
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string, string, int, int) ([]*entity.Person, int, error)); ok {
		return rf(email, phone, firstName, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(string, string, string, int, int) []*entity.Person); ok {
		r0 = rf(email, phone, firstName, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.Person)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string, string, int, int) int); ok {
		r1 = rf(email, phone, firstName, limit, offset)
	} else {
		r1 = ret.Get(1).(int)
	}

	if rf, ok := ret.Get(2).(func(string, string, string, int, int) error); ok {
		r2 = rf(email, phone, firstName, limit, offset)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetByEmail provides a mock function with given fields: email
func (_m *PersonLogic) GetByEmail(email string) (*entity.Person, error) {
	ret := _m.Called(email)

	var r0 *entity.Person
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*entity.Person, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *entity.Person); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Person)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: id
func (_m *PersonLogic) GetByID(id int) (*entity.Person, error) {
	ret := _m.Called(id)

	var r0 *entity.Person
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (*entity.Person, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int) *entity.Person); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Person)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: id, req
func (_m *PersonLogic) Update(id int, req *entity.Person) (*entity.Person, error) {
	ret := _m.Called(id, req)

	var r0 *entity.Person
	var r1 error
	if rf, ok := ret.Get(0).(func(int, *entity.Person) (*entity.Person, error)); ok {
		return rf(id, req)
	}
	if rf, ok := ret.Get(0).(func(int, *entity.Person) *entity.Person); ok {
		r0 = rf(id, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Person)
		}
	}

	if rf, ok := ret.Get(1).(func(int, *entity.Person) error); ok {
		r1 = rf(id, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPersonLogic creates a new instance of PersonLogic. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPersonLogic(t interface {
	mock.TestingT
	Cleanup(func())
}) *PersonLogic {
	mock := &PersonLogic{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
