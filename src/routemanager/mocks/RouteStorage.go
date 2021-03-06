// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import domain "github.com/JaneKetko/Buses/src/domain"
import mock "github.com/stretchr/testify/mock"

// RouteStorage is an autogenerated mock type for the RouteStorage type
type RouteStorage struct {
	mock.Mock
}

// AddRoute provides a mock function with given fields: _a0
func (_m *RouteStorage) AddRoute(_a0 *domain.Route) (int, error) {
	ret := _m.Called(_a0)

	var r0 int
	if rf, ok := ret.Get(0).(func(*domain.Route) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*domain.Route) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteRow provides a mock function with given fields: id
func (_m *RouteStorage) DeleteRow(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllData provides a mock function with given fields:
func (_m *RouteStorage) GetAllData() ([]domain.Route, error) {
	ret := _m.Called()

	var r0 []domain.Route
	if rf, ok := ret.Get(0).(func() []domain.Route); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Route)
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

// RouteByID provides a mock function with given fields: id
func (_m *RouteStorage) RouteByID(id int) (*domain.Route, error) {
	ret := _m.Called(id)

	var r0 *domain.Route
	if rf, ok := ret.Get(0).(func(int) *domain.Route); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Route)
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

// RoutesByEndPoint provides a mock function with given fields: point
func (_m *RouteStorage) RoutesByEndPoint(point string) ([]domain.Route, error) {
	ret := _m.Called(point)

	var r0 []domain.Route
	if rf, ok := ret.Get(0).(func(string) []domain.Route); ok {
		r0 = rf(point)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Route)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(point)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
