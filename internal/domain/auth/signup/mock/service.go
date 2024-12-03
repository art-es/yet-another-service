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

// MockhashGenerator is a mock of hashGenerator interface.
type MockhashGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockhashGeneratorMockRecorder
	isgomock struct{}
}

// MockhashGeneratorMockRecorder is the mock recorder for MockhashGenerator.
type MockhashGeneratorMockRecorder struct {
	mock *MockhashGenerator
}

// NewMockhashGenerator creates a new mock instance.
func NewMockhashGenerator(ctrl *gomock.Controller) *MockhashGenerator {
	mock := &MockhashGenerator{ctrl: ctrl}
	mock.recorder = &MockhashGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockhashGenerator) EXPECT() *MockhashGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockhashGenerator) Generate(str string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", str)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockhashGeneratorMockRecorder) Generate(str any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockhashGenerator)(nil).Generate), str)
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

// Add mocks base method.
func (m *MockuserRepository) Add(ctx context.Context, tx transaction.Transaction, user *auth.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, tx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockuserRepositoryMockRecorder) Add(ctx, tx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockuserRepository)(nil).Add), ctx, tx, user)
}

// EmailExists mocks base method.
func (m *MockuserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EmailExists", ctx, email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EmailExists indicates an expected call of EmailExists.
func (mr *MockuserRepositoryMockRecorder) EmailExists(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EmailExists", reflect.TypeOf((*MockuserRepository)(nil).EmailExists), ctx, email)
}

// MockactivationCreator is a mock of activationCreator interface.
type MockactivationCreator struct {
	ctrl     *gomock.Controller
	recorder *MockactivationCreatorMockRecorder
	isgomock struct{}
}

// MockactivationCreatorMockRecorder is the mock recorder for MockactivationCreator.
type MockactivationCreatorMockRecorder struct {
	mock *MockactivationCreator
}

// NewMockactivationCreator creates a new mock instance.
func NewMockactivationCreator(ctrl *gomock.Controller) *MockactivationCreator {
	mock := &MockactivationCreator{ctrl: ctrl}
	mock.recorder = &MockactivationCreatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockactivationCreator) EXPECT() *MockactivationCreatorMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockactivationCreator) Create(ctx context.Context, tx transaction.Transaction, userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, tx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockactivationCreatorMockRecorder) Create(ctx, tx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockactivationCreator)(nil).Create), ctx, tx, userID)
}