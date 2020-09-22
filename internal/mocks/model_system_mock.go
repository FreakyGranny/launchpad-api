// Code generated by MockGen. DO NOT EDIT.
// Source: system.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "github.com/FreakyGranny/launchpad-api/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSystemImpl is a mock of SystemImpl interface
type MockSystemImpl struct {
	ctrl     *gomock.Controller
	recorder *MockSystemImplMockRecorder
}

// MockSystemImplMockRecorder is the mock recorder for MockSystemImpl
type MockSystemImplMockRecorder struct {
	mock *MockSystemImpl
}

// NewMockSystemImpl creates a new mock instance
func NewMockSystemImpl(ctrl *gomock.Controller) *MockSystemImpl {
	mock := &MockSystemImpl{ctrl: ctrl}
	mock.recorder = &MockSystemImplMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSystemImpl) EXPECT() *MockSystemImplMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockSystemImpl) Get() (*models.System, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(*models.System)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockSystemImplMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSystemImpl)(nil).Get))
}

// Update mocks base method
func (m *MockSystemImpl) Update(s *models.System) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", s)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockSystemImplMockRecorder) Update(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSystemImpl)(nil).Update), s)
}