// Code generated by MockGen. DO NOT EDIT.
// Source: queue.go

// Package mock_redisqueue is a generated GoMock package.
package mock_redisqueue

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockQueue is a mock of Queue interface
type MockQueue struct {
	ctrl     *gomock.Controller
	recorder *MockQueueMockRecorder
}

// MockQueueMockRecorder is the mock recorder for MockQueue
type MockQueueMockRecorder struct {
	mock *MockQueue
}

// NewMockQueue creates a new mock instance
func NewMockQueue(ctrl *gomock.Controller) *MockQueue {
	mock := &MockQueue{ctrl: ctrl}
	mock.recorder = &MockQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQueue) EXPECT() *MockQueueMockRecorder {
	return m.recorder
}

// Push mocks base method
func (m *MockQueue) Push(key, val string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Push", key, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push
func (mr *MockQueueMockRecorder) Push(key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockQueue)(nil).Push), key, val)
}

// Peek mocks base method
func (m *MockQueue) Peek(key string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peek", key)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Peek indicates an expected call of Peek
func (mr *MockQueueMockRecorder) Peek(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peek", reflect.TypeOf((*MockQueue)(nil).Peek), key)
}

// Remove mocks base method
func (m *MockQueue) Remove(key string, count int64, val interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", key, count, val)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockQueueMockRecorder) Remove(key, count, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockQueue)(nil).Remove), key, count, val)
}

// PopPush mocks base method
func (m *MockQueue) PopPush(source, destination string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PopPush", source, destination)
	ret0, _ := ret[0].(error)
	return ret0
}

// PopPush indicates an expected call of PopPush
func (mr *MockQueueMockRecorder) PopPush(source, destination interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopPush", reflect.TypeOf((*MockQueue)(nil).PopPush), source, destination)
}