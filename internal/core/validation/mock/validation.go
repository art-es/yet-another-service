// Code generated by MockGen. DO NOT EDIT.
// Source: validation.go
//
// Generated by this command:
//
//	mockgen -source=validation.go -destination=mock/validation.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockValidator is a mock of Validator interface.
type MockValidator struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorMockRecorder
	isgomock struct{}
}

// MockValidatorMockRecorder is the mock recorder for MockValidator.
type MockValidatorMockRecorder struct {
	mock *MockValidator
}

// NewMockValidator creates a new mock instance.
func NewMockValidator(ctrl *gomock.Controller) *MockValidator {
	mock := &MockValidator{ctrl: ctrl}
	mock.recorder = &MockValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockValidator) EXPECT() *MockValidatorMockRecorder {
	return m.recorder
}

// Struct mocks base method.
func (m *MockValidator) Struct(s any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Struct", s)
	ret0, _ := ret[0].(error)
	return ret0
}

// Struct indicates an expected call of Struct.
func (mr *MockValidatorMockRecorder) Struct(s any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Struct", reflect.TypeOf((*MockValidator)(nil).Struct), s)
}

// Var mocks base method.
func (m *MockValidator) Var(field any, tag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Var", field, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// Var indicates an expected call of Var.
func (mr *MockValidatorMockRecorder) Var(field, tag any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Var", reflect.TypeOf((*MockValidator)(nil).Var), field, tag)
}
