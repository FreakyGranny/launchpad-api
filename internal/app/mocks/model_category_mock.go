// Code generated by MockGen. DO NOT EDIT.
// Source: category.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "github.com/FreakyGranny/launchpad-api/internal/app/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCategoryImpl is a mock of CategoryImpl interface
type MockCategoryImpl struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryImplMockRecorder
}

// MockCategoryImplMockRecorder is the mock recorder for MockCategoryImpl
type MockCategoryImplMockRecorder struct {
	mock *MockCategoryImpl
}

// NewMockCategoryImpl creates a new mock instance
func NewMockCategoryImpl(ctrl *gomock.Controller) *MockCategoryImpl {
	mock := &MockCategoryImpl{ctrl: ctrl}
	mock.recorder = &MockCategoryImplMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCategoryImpl) EXPECT() *MockCategoryImplMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockCategoryImpl) Get(id int) (*models.Category, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*models.Category)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockCategoryImplMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCategoryImpl)(nil).Get), id)
}

// GetAll mocks base method
func (m *MockCategoryImpl) GetAll() ([]models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockCategoryImplMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockCategoryImpl)(nil).GetAll))
}
