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
	models "github.com/art-es/yet-another-service/internal/domain/shared/models"
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

// Exists mocks base method.
func (m *MockuserRepository) Exists(ctx context.Context, email string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", ctx, email)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockuserRepositoryMockRecorder) Exists(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockuserRepository)(nil).Exists), ctx, email)
}

// Save mocks base method.
func (m *MockuserRepository) Save(ctx context.Context, tx transaction.Transaction, user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, tx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockuserRepositoryMockRecorder) Save(ctx, tx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockuserRepository)(nil).Save), ctx, tx, user)
}

// MockactivationService is a mock of activationService interface.
type MockactivationService struct {
	ctrl     *gomock.Controller
	recorder *MockactivationServiceMockRecorder
	isgomock struct{}
}

// MockactivationServiceMockRecorder is the mock recorder for MockactivationService.
type MockactivationServiceMockRecorder struct {
	mock *MockactivationService
}

// NewMockactivationService creates a new mock instance.
func NewMockactivationService(ctrl *gomock.Controller) *MockactivationService {
	mock := &MockactivationService{ctrl: ctrl}
	mock.recorder = &MockactivationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockactivationService) EXPECT() *MockactivationServiceMockRecorder {
	return m.recorder
}

// CreateActivation mocks base method.
func (m *MockactivationService) CreateActivation(ctx context.Context, tx transaction.Transaction, user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActivation", ctx, tx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateActivation indicates an expected call of CreateActivation.
func (mr *MockactivationServiceMockRecorder) CreateActivation(ctx, tx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActivation", reflect.TypeOf((*MockactivationService)(nil).CreateActivation), ctx, tx, user)
}
