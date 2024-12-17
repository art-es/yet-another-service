// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mock/service.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MocktokenService is a mock of tokenService interface.
type MocktokenService struct {
	ctrl     *gomock.Controller
	recorder *MocktokenServiceMockRecorder
	isgomock struct{}
}

// MocktokenServiceMockRecorder is the mock recorder for MocktokenService.
type MocktokenServiceMockRecorder struct {
	mock *MocktokenService
}

// NewMocktokenService creates a new mock instance.
func NewMocktokenService(ctrl *gomock.Controller) *MocktokenService {
	mock := &MocktokenService{ctrl: ctrl}
	mock.recorder = &MocktokenServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocktokenService) EXPECT() *MocktokenServiceMockRecorder {
	return m.recorder
}

// Invalidate mocks base method.
func (m *MocktokenService) Invalidate(ctx context.Context, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Invalidate", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Invalidate indicates an expected call of Invalidate.
func (mr *MocktokenServiceMockRecorder) Invalidate(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invalidate", reflect.TypeOf((*MocktokenService)(nil).Invalidate), ctx, token)
}