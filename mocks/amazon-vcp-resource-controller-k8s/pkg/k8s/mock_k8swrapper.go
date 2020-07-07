// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aws/amazon-vpc-resource-controller-k8s/pkg/k8s (interfaces: K8sWrapper)

// Package mock_k8s is a generated GoMock package.
package mock_k8s

import (
	v1alpha1 "github.com/aws/amazon-vpc-cni-k8s/pkg/apis/crd/v1alpha1"
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	reflect "reflect"
)

// MockK8sWrapper is a mock of K8sWrapper interface
type MockK8sWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockK8sWrapperMockRecorder
}

// MockK8sWrapperMockRecorder is the mock recorder for MockK8sWrapper
type MockK8sWrapperMockRecorder struct {
	mock *MockK8sWrapper
}

// NewMockK8sWrapper creates a new mock instance
func NewMockK8sWrapper(ctrl *gomock.Controller) *MockK8sWrapper {
	mock := &MockK8sWrapper{ctrl: ctrl}
	mock.recorder = &MockK8sWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockK8sWrapper) EXPECT() *MockK8sWrapperMockRecorder {
	return m.recorder
}

// AdvertiseCapacityIfNotSet mocks base method
func (m *MockK8sWrapper) AdvertiseCapacityIfNotSet(arg0, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdvertiseCapacityIfNotSet", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AdvertiseCapacityIfNotSet indicates an expected call of AdvertiseCapacityIfNotSet
func (mr *MockK8sWrapperMockRecorder) AdvertiseCapacityIfNotSet(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdvertiseCapacityIfNotSet", reflect.TypeOf((*MockK8sWrapper)(nil).AdvertiseCapacityIfNotSet), arg0, arg1, arg2)
}

// AnnotatePod mocks base method
func (m *MockK8sWrapper) AnnotatePod(arg0, arg1, arg2, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AnnotatePod", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AnnotatePod indicates an expected call of AnnotatePod
func (mr *MockK8sWrapperMockRecorder) AnnotatePod(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AnnotatePod", reflect.TypeOf((*MockK8sWrapper)(nil).AnnotatePod), arg0, arg1, arg2, arg3)
}

// GetENIConfig mocks base method
func (m *MockK8sWrapper) GetENIConfig(arg0 string) (*v1alpha1.ENIConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetENIConfig", arg0)
	ret0, _ := ret[0].(*v1alpha1.ENIConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetENIConfig indicates an expected call of GetENIConfig
func (mr *MockK8sWrapperMockRecorder) GetENIConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetENIConfig", reflect.TypeOf((*MockK8sWrapper)(nil).GetENIConfig), arg0)
}

// GetPod mocks base method
func (m *MockK8sWrapper) GetPod(arg0, arg1 string) (*v1.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPod", arg0, arg1)
	ret0, _ := ret[0].(*v1.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPod indicates an expected call of GetPod
func (mr *MockK8sWrapperMockRecorder) GetPod(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPod", reflect.TypeOf((*MockK8sWrapper)(nil).GetPod), arg0, arg1)
}

// GetPodFromAPIServer mocks base method
func (m *MockK8sWrapper) GetPodFromAPIServer(arg0, arg1 string) (*v1.Pod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPodFromAPIServer", arg0, arg1)
	ret0, _ := ret[0].(*v1.Pod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPodFromAPIServer indicates an expected call of GetPodFromAPIServer
func (mr *MockK8sWrapperMockRecorder) GetPodFromAPIServer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPodFromAPIServer", reflect.TypeOf((*MockK8sWrapper)(nil).GetPodFromAPIServer), arg0, arg1)
}

// ListPods mocks base method
func (m *MockK8sWrapper) ListPods(arg0 string) (*v1.PodList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPods", arg0)
	ret0, _ := ret[0].(*v1.PodList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPods indicates an expected call of ListPods
func (mr *MockK8sWrapperMockRecorder) ListPods(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPods", reflect.TypeOf((*MockK8sWrapper)(nil).ListPods), arg0)
}
