// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go
//
// Generated by this command:
//
//	mockgen -source=contract.go -destination=mock/contract.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockMailer is a mock of Mailer interface.
type MockMailer struct {
	ctrl     *gomock.Controller
	recorder *MockMailerMockRecorder
	isgomock struct{}
}

// MockMailerMockRecorder is the mock recorder for MockMailer.
type MockMailerMockRecorder struct {
	mock *MockMailer
}

// NewMockMailer creates a new mock instance.
func NewMockMailer(ctrl *gomock.Controller) *MockMailer {
	mock := &MockMailer{ctrl: ctrl}
	mock.recorder = &MockMailerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMailer) EXPECT() *MockMailerMockRecorder {
	return m.recorder
}

// MailTo mocks base method.
func (m *MockMailer) MailTo(emailAddress, subject, content string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MailTo", emailAddress, subject, content)
	ret0, _ := ret[0].(error)
	return ret0
}

// MailTo indicates an expected call of MailTo.
func (mr *MockMailerMockRecorder) MailTo(emailAddress, subject, content any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MailTo", reflect.TypeOf((*MockMailer)(nil).MailTo), emailAddress, subject, content)
}
