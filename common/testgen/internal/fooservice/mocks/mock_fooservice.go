// Code generated by MockGen. DO NOT EDIT.
// Source: ./fooservice.go

// Package mock_fooservice is a generated GoMock package.
package mock_fooservice

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockExternalService is a mock of ExternalService interface.
type MockExternalService struct {
	ctrl     *gomock.Controller
	recorder *MockExternalServiceMockRecorder
}

// MockExternalServiceMockRecorder is the mock recorder for MockExternalService.
type MockExternalServiceMockRecorder struct {
	mock *MockExternalService
}

// NewMockExternalService creates a new mock instance.
func NewMockExternalService(ctrl *gomock.Controller) *MockExternalService {
	mock := &MockExternalService{ctrl: ctrl}
	mock.recorder = &MockExternalServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExternalService) EXPECT() *MockExternalServiceMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockExternalService) Call() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call")
	ret0, _ := ret[0].(error)
	return ret0
}

// Call indicates an expected call of Call.
func (mr *MockExternalServiceMockRecorder) Call() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockExternalService)(nil).Call))
}
