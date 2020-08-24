// Code generated by MockGen. DO NOT EDIT.
// Source: donation.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "github.com/FreakyGranny/launchpad-api/internal/app/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDonationImpl is a mock of DonationImpl interface
type MockDonationImpl struct {
	ctrl     *gomock.Controller
	recorder *MockDonationImplMockRecorder
}

// MockDonationImplMockRecorder is the mock recorder for MockDonationImpl
type MockDonationImplMockRecorder struct {
	mock *MockDonationImpl
}

// NewMockDonationImpl creates a new mock instance
func NewMockDonationImpl(ctrl *gomock.Controller) *MockDonationImpl {
	mock := &MockDonationImpl{ctrl: ctrl}
	mock.recorder = &MockDonationImplMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDonationImpl) EXPECT() *MockDonationImplMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockDonationImpl) Get(id int) (*models.Donation, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(*models.Donation)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockDonationImplMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDonationImpl)(nil).Get), id)
}

// GetAllByUser mocks base method
func (m *MockDonationImpl) GetAllByUser(id int) ([]models.Donation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUser", id)
	ret0, _ := ret[0].([]models.Donation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByUser indicates an expected call of GetAllByUser
func (mr *MockDonationImplMockRecorder) GetAllByUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUser", reflect.TypeOf((*MockDonationImpl)(nil).GetAllByUser), id)
}

// GetAllByProject mocks base method
func (m *MockDonationImpl) GetAllByProject(id int) ([]models.Donation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByProject", id)
	ret0, _ := ret[0].([]models.Donation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByProject indicates an expected call of GetAllByProject
func (mr *MockDonationImplMockRecorder) GetAllByProject(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByProject", reflect.TypeOf((*MockDonationImpl)(nil).GetAllByProject), id)
}

// Create mocks base method
func (m *MockDonationImpl) Create(d *models.Donation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", d)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockDonationImplMockRecorder) Create(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDonationImpl)(nil).Create), d)
}

// Update mocks base method
func (m *MockDonationImpl) Update(d *models.Donation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", d)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockDonationImplMockRecorder) Update(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDonationImpl)(nil).Update), d)
}

// Delete mocks base method
func (m *MockDonationImpl) Delete(d *models.Donation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", d)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockDonationImplMockRecorder) Delete(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDonationImpl)(nil).Delete), d)
}
