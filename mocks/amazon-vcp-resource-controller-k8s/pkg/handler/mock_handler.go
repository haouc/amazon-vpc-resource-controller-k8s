// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aws/amazon-vpc-resource-controller-k8s/pkg/handler (interfaces: Handler)

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	reflect "reflect"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// MockHandler is a mock of Handler interface
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// CanHandle mocks base method
func (m *MockHandler) CanHandle(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CanHandle", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CanHandle indicates an expected call of CanHandle
func (mr *MockHandlerMockRecorder) CanHandle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CanHandle", reflect.TypeOf((*MockHandler)(nil).CanHandle), arg0)
}

// HandleCreate mocks base method
func (m *MockHandler) HandleCreate(arg0 string, arg1 int, arg2 *v1.Pod) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleCreate", arg0, arg1, arg2)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleCreate indicates an expected call of HandleCreate
func (mr *MockHandlerMockRecorder) HandleCreate(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleCreate", reflect.TypeOf((*MockHandler)(nil).HandleCreate), arg0, arg1, arg2)
}

// HandleDelete mocks base method
func (m *MockHandler) HandleDelete(arg0 string, arg1 *v1.Pod) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleDelete", arg0, arg1)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HandleDelete indicates an expected call of HandleDelete
func (mr *MockHandlerMockRecorder) HandleDelete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleDelete", reflect.TypeOf((*MockHandler)(nil).HandleDelete), arg0, arg1)
}
