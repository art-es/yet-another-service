// Code generated by MockGen. DO NOT EDIT.
// Source: retrier.go
//
// Generated by this command:
//
//	mockgen -source=retrier.go -destination=mock/retrier.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRetrier is a mock of Retrier interface.
type MockRetrier struct {
	ctrl     *gomock.Controller
	recorder *MockRetrierMockRecorder
	isgomock struct{}
}

// MockRetrierMockRecorder is the mock recorder for MockRetrier.
type MockRetrierMockRecorder struct {
	mock *MockRetrier
}

// NewMockRetrier creates a new mock instance.
func NewMockRetrier(ctrl *gomock.Controller) *MockRetrier {
	mock := &MockRetrier{ctrl: ctrl}
	mock.recorder = &MockRetrierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRetrier) EXPECT() *MockRetrierMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockRetrier) Process(f func() error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", f)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockRetrierMockRecorder) Process(f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockRetrier)(nil).Process), f)
}
