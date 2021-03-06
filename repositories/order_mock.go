// Code generated by MockGen. DO NOT EDIT.
// Source: order.go

// Package repositories is a generated GoMock package.
package repositories

import (
	gomock "github.com/golang/mock/gomock"
	models "github.com/netology/dao-pattern/models"
	reflect "reflect"
)

// MockOrderRepository is a mock of OrderRepository interface
type MockOrderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryMockRecorder
}

// MockOrderRepositoryMockRecorder is the mock recorder for MockOrderRepository
type MockOrderRepositoryMockRecorder struct {
	mock *MockOrderRepository
}

// NewMockOrderRepository creates a new mock instance
func NewMockOrderRepository(ctrl *gomock.Controller) *MockOrderRepository {
	mock := &MockOrderRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOrderRepository) EXPECT() *MockOrderRepositoryMockRecorder {
	return m.recorder
}

// GetByID mocks base method
func (m *MockOrderRepository) GetByID(orderID models.OrderID) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", orderID)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockOrderRepositoryMockRecorder) GetByID(orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockOrderRepository)(nil).GetByID), orderID)
}

// Save mocks base method
func (m *MockOrderRepository) Save(order *models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockOrderRepositoryMockRecorder) Save(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockOrderRepository)(nil).Save), order)
}
