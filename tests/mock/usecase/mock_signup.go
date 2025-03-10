// Code generated by MockGen. DO NOT EDIT.
// Source: signup.go
//
// Generated by this command:
//
//	mockgen -source=signup.go -destination=../tests/mock/usecase/mock_signup.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSignUpUseCase is a mock of SignUpUseCase interface.
type MockSignUpUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSignUpUseCaseMockRecorder
	isgomock struct{}
}

// MockSignUpUseCaseMockRecorder is the mock recorder for MockSignUpUseCase.
type MockSignUpUseCaseMockRecorder struct {
	mock *MockSignUpUseCase
}

// NewMockSignUpUseCase creates a new mock instance.
func NewMockSignUpUseCase(ctrl *gomock.Controller) *MockSignUpUseCase {
	mock := &MockSignUpUseCase{ctrl: ctrl}
	mock.recorder = &MockSignUpUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSignUpUseCase) EXPECT() *MockSignUpUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockSignUpUseCase) Execute(ctx context.Context, email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute.
func (mr *MockSignUpUseCaseMockRecorder) Execute(ctx, email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockSignUpUseCase)(nil).Execute), ctx, email, password)
}
