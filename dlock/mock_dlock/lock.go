// Code generated by MockGen. DO NOT EDIT.
// Source: lock.go

// Package mock_dlock is a generated GoMock package.
package mock_dlock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockLock is a mock of Lock interface
type MockLock struct {
	ctrl     *gomock.Controller
	recorder *MockLockMockRecorder
}

// MockLockMockRecorder is the mock recorder for MockLock
type MockLockMockRecorder struct {
	mock *MockLock
}

// NewMockLock creates a new mock instance
func NewMockLock(ctrl *gomock.Controller) *MockLock {
	mock := &MockLock{ctrl: ctrl}
	mock.recorder = &MockLockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLock) EXPECT() *MockLockMockRecorder {
	return m.recorder
}

// Lock mocks base method
func (m *MockLock) Lock(id string, expiry time.Duration) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lock", id, expiry)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lock indicates an expected call of Lock
func (mr *MockLockMockRecorder) Lock(id, expiry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lock", reflect.TypeOf((*MockLock)(nil).Lock), id, expiry)
}

// Unlock mocks base method
func (m *MockLock) Unlock(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unlock", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unlock indicates an expected call of Unlock
func (mr *MockLockMockRecorder) Unlock(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unlock", reflect.TypeOf((*MockLock)(nil).Unlock), id)
}

// IsLocked mocks base method
func (m *MockLock) IsLocked(id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsLocked", id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsLocked indicates an expected call of IsLocked
func (mr *MockLockMockRecorder) IsLocked(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsLocked", reflect.TypeOf((*MockLock)(nil).IsLocked), id)
}
