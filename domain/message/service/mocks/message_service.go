// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/santileira/fuck-off-as-a-service/domain/message/domain"
	mock "github.com/stretchr/testify/mock"
)

// MessageService is an autogenerated mock type for the MessageService type
type MessageService struct {
	mock.Mock
}

// GetMessage provides a mock function with given fields: userID
func (_m *MessageService) GetMessage(userID string) (*domain.Response, error) {
	ret := _m.Called(userID)

	var r0 *domain.Response
	if rf, ok := ret.Get(0).(func(string) *domain.Response); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
