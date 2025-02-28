// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go
//
// Generated by this command:
//
//	mockgen -source=handler.go -destination=mock/handler.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockrecoveryService is a mock of recoveryService interface.
type MockrecoveryService struct {
	ctrl     *gomock.Controller
	recorder *MockrecoveryServiceMockRecorder
	isgomock struct{}
}

// MockrecoveryServiceMockRecorder is the mock recorder for MockrecoveryService.
type MockrecoveryServiceMockRecorder struct {
	mock *MockrecoveryService
}

// NewMockrecoveryService creates a new mock instance.
func NewMockrecoveryService(ctrl *gomock.Controller) *MockrecoveryService {
	mock := &MockrecoveryService{ctrl: ctrl}
	mock.recorder = &MockrecoveryServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockrecoveryService) EXPECT() *MockrecoveryServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockrecoveryService) Create(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockrecoveryServiceMockRecorder) Create(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockrecoveryService)(nil).Create), ctx, email)
}
