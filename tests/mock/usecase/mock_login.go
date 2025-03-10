// Code generated by MockGen. DO NOT EDIT.
// Source: login.go
//
// Generated by this command:
//
//	mockgen -source=login.go -destination=../tests/mock/usecase/mock_login.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLoginUseCase is a mock of LoginUseCase interface.
type MockLoginUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockLoginUseCaseMockRecorder
	isgomock struct{}
}

// MockLoginUseCaseMockRecorder is the mock recorder for MockLoginUseCase.
type MockLoginUseCaseMockRecorder struct {
	mock *MockLoginUseCase
}

// NewMockLoginUseCase creates a new mock instance.
func NewMockLoginUseCase(ctrl *gomock.Controller) *MockLoginUseCase {
	mock := &MockLoginUseCase{ctrl: ctrl}
	mock.recorder = &MockLoginUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoginUseCase) EXPECT() *MockLoginUseCaseMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockLoginUseCase) Execute(ctx context.Context, email, password string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", ctx, email, password)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockLoginUseCaseMockRecorder) Execute(ctx, email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockLoginUseCase)(nil).Execute), ctx, email, password)
}
