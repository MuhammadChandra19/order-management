// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	product "github.com/MuhammadChandra19/order-management/internal/product"
	gomock "github.com/golang/mock/gomock"
)

// MockProductRepositoryInterface is a mock of ProductRepositoryInterface interface.
type MockProductRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockProductRepositoryInterfaceMockRecorder
}

// MockProductRepositoryInterfaceMockRecorder is the mock recorder for MockProductRepositoryInterface.
type MockProductRepositoryInterfaceMockRecorder struct {
	mock *MockProductRepositoryInterface
}

// NewMockProductRepositoryInterface creates a new mock instance.
func NewMockProductRepositoryInterface(ctrl *gomock.Controller) *MockProductRepositoryInterface {
	mock := &MockProductRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockProductRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductRepositoryInterface) EXPECT() *MockProductRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetProductSalesStats mocks base method.
func (m *MockProductRepositoryInterface) GetProductSalesStats() ([]product.ProductSalesStats, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductSalesStats")
	ret0, _ := ret[0].([]product.ProductSalesStats)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductSalesStats indicates an expected call of GetProductSalesStats.
func (mr *MockProductRepositoryInterfaceMockRecorder) GetProductSalesStats() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductSalesStats", reflect.TypeOf((*MockProductRepositoryInterface)(nil).GetProductSalesStats))
}