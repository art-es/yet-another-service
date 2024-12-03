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

	transaction "github.com/art-es/yet-another-service/internal/core/transaction"
	auth "github.com/art-es/yet-another-service/internal/domain/auth"
	gomock "go.uber.org/mock/gomock"
)

// MockactivationRepository is a mock of activationRepository interface.
type MockactivationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockactivationRepositoryMockRecorder
	isgomock struct{}
}

// MockactivationRepositoryMockRecorder is the mock recorder for MockactivationRepository.
type MockactivationRepositoryMockRecorder struct {
	mock *MockactivationRepository
}

// NewMockactivationRepository creates a new mock instance.
func NewMockactivationRepository(ctrl *gomock.Controller) *MockactivationRepository {
	mock := &MockactivationRepository{ctrl: ctrl}
	mock.recorder = &MockactivationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockactivationRepository) EXPECT() *MockactivationRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockactivationRepository) Delete(ctx context.Context, tx transaction.Transaction, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, tx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockactivationRepositoryMockRecorder) Delete(ctx, tx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockactivationRepository)(nil).Delete), ctx, tx, token)
}

// FindByToken mocks base method.
func (m *MockactivationRepository) FindByToken(ctx context.Context, token string) (*auth.Activation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByToken", ctx, token)
	ret0, _ := ret[0].(*auth.Activation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByToken indicates an expected call of FindByToken.
func (mr *MockactivationRepositoryMockRecorder) FindByToken(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByToken", reflect.TypeOf((*MockactivationRepository)(nil).FindByToken), ctx, token)
}

// MockuserRepository is a mock of userRepository interface.
type MockuserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockuserRepositoryMockRecorder
	isgomock struct{}
}

// MockuserRepositoryMockRecorder is the mock recorder for MockuserRepository.
type MockuserRepositoryMockRecorder struct {
	mock *MockuserRepository
}

// NewMockuserRepository creates a new mock instance.
func NewMockuserRepository(ctrl *gomock.Controller) *MockuserRepository {
	mock := &MockuserRepository{ctrl: ctrl}
	mock.recorder = &MockuserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserRepository) EXPECT() *MockuserRepositoryMockRecorder {
	return m.recorder
}

// Activate mocks base method.
func (m *MockuserRepository) Activate(ctx context.Context, tx transaction.Transaction, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Activate", ctx, tx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Activate indicates an expected call of Activate.
func (mr *MockuserRepositoryMockRecorder) Activate(ctx, tx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Activate", reflect.TypeOf((*MockuserRepository)(nil).Activate), ctx, tx, userID)
}