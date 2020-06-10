// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/aws/amazon-vpc-resource-controller-k8s/pkg/aws/ec2 (interfaces: EC2MetadataClient)

// Package mock_ec2 is a generated GoMock package.
package mock_ec2

import (
	ec2metadata "github.com/aws/aws-sdk-go/aws/ec2metadata"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockEC2MetadataClient is a mock of EC2MetadataClient interface
type MockEC2MetadataClient struct {
	ctrl     *gomock.Controller
	recorder *MockEC2MetadataClientMockRecorder
}

// MockEC2MetadataClientMockRecorder is the mock recorder for MockEC2MetadataClient
type MockEC2MetadataClientMockRecorder struct {
	mock *MockEC2MetadataClient
}

// NewMockEC2MetadataClient creates a new mock instance
func NewMockEC2MetadataClient(ctrl *gomock.Controller) *MockEC2MetadataClient {
	mock := &MockEC2MetadataClient{ctrl: ctrl}
	mock.recorder = &MockEC2MetadataClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEC2MetadataClient) EXPECT() *MockEC2MetadataClientMockRecorder {
	return m.recorder
}

// GetInstanceIdentityDocument mocks base method
func (m *MockEC2MetadataClient) GetInstanceIdentityDocument() (ec2metadata.EC2InstanceIdentityDocument, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstanceIdentityDocument")
	ret0, _ := ret[0].(ec2metadata.EC2InstanceIdentityDocument)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstanceIdentityDocument indicates an expected call of GetInstanceIdentityDocument
func (mr *MockEC2MetadataClientMockRecorder) GetInstanceIdentityDocument() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstanceIdentityDocument", reflect.TypeOf((*MockEC2MetadataClient)(nil).GetInstanceIdentityDocument))
}

// Region mocks base method
func (m *MockEC2MetadataClient) Region() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Region")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Region indicates an expected call of Region
func (mr *MockEC2MetadataClientMockRecorder) Region() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Region", reflect.TypeOf((*MockEC2MetadataClient)(nil).Region))
}
