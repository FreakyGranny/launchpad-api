// Code generated by MockGen. DO NOT EDIT.
// Source: project.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "github.com/FreakyGranny/launchpad-api/internal/app/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockProjectImpl is a mock of ProjectImpl interface
type MockProjectImpl struct {
	ctrl     *gomock.Controller
	recorder *MockProjectImplMockRecorder
}

// MockProjectImplMockRecorder is the mock recorder for MockProjectImpl
type MockProjectImplMockRecorder struct {
	mock *MockProjectImpl
}

// NewMockProjectImpl creates a new mock instance
func NewMockProjectImpl(ctrl *gomock.Controller) *MockProjectImpl {
	mock := &MockProjectImpl{ctrl: ctrl}
	mock.recorder = &MockProjectImplMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockProjectImpl) EXPECT() *MockProjectImplMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockProjectImpl) Get(id int) (*models.Project, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*models.Project)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockProjectImplMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockProjectImpl)(nil).Get), id)
}