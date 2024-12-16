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
	reflect "reflect"

	auth "github.com/art-es/yet-another-service/internal/app/auth"
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

// Generate mocks base method.
func (m *MocktokenService) Generate(claims *auth.TokenClaims) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", claims)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MocktokenServiceMockRecorder) Generate(claims any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MocktokenService)(nil).Generate), claims)
}

// Parse mocks base method.
func (m *MocktokenService) Parse(token string) (*auth.TokenClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", token)
	ret0, _ := ret[0].(*auth.TokenClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse.
func (mr *MocktokenServiceMockRecorder) Parse(token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MocktokenService)(nil).Parse), token)
}