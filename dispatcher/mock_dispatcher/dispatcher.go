// Code generated by MockGen. DO NOT EDIT.
// Source: dispatcher.go

// Package mock_dispatcher is a generated GoMock package.
package mock_dispatcher

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDispatcher is a mock of Dispatcher interface
type MockDispatcher struct {
	ctrl     *gomock.Controller
	recorder *MockDispatcherMockRecorder
}

// MockDispatcherMockRecorder is the mock recorder for MockDispatcher
type MockDispatcherMockRecorder struct {
	mock *MockDispatcher
}

// NewMockDispatcher creates a new mock instance
func NewMockDispatcher(ctrl *gomock.Controller) *MockDispatcher {
	mock := &MockDispatcher{ctrl: ctrl}
	mock.recorder = &MockDispatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDispatcher) EXPECT() *MockDispatcherMockRecorder {
	return m.recorder
}

// Post mocks base method
func (m *MockDispatcher) Post(payload, jobType string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Post", payload, jobType)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Post indicates an expected call of Post
func (mr *MockDispatcherMockRecorder) Post(payload, jobType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Post", reflect.TypeOf((*MockDispatcher)(nil).Post), payload, jobType)
}
