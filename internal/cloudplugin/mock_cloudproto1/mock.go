// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/opentofu/opentofu/internal/cloudplugin/cloudproto1 (interfaces: CommandServiceClient,CommandService_ExecuteClient)

// Package mock_cloudproto1 is a generated GoMock package.
package mock_cloudproto1

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cloudproto1 "github.com/opentofu/opentofu/internal/cloudplugin/cloudproto1"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockCommandServiceClient is a mock of CommandServiceClient interface.
type MockCommandServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockCommandServiceClientMockRecorder
}

// MockCommandServiceClientMockRecorder is the mock recorder for MockCommandServiceClient.
type MockCommandServiceClientMockRecorder struct {
	mock *MockCommandServiceClient
}

// NewMockCommandServiceClient creates a new mock instance.
func NewMockCommandServiceClient(ctrl *gomock.Controller) *MockCommandServiceClient {
	mock := &MockCommandServiceClient{ctrl: ctrl}
	mock.recorder = &MockCommandServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandServiceClient) EXPECT() *MockCommandServiceClientMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockCommandServiceClient) Execute(arg0 context.Context, arg1 *cloudproto1.CommandRequest, arg2 ...grpc.CallOption) (cloudproto1.CommandService_ExecuteClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(cloudproto1.CommandService_ExecuteClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockCommandServiceClientMockRecorder) Execute(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCommandServiceClient)(nil).Execute), varargs...)
}

// MockCommandService_ExecuteClient is a mock of CommandService_ExecuteClient interface.
type MockCommandService_ExecuteClient struct {
	ctrl     *gomock.Controller
	recorder *MockCommandService_ExecuteClientMockRecorder
}

// MockCommandService_ExecuteClientMockRecorder is the mock recorder for MockCommandService_ExecuteClient.
type MockCommandService_ExecuteClientMockRecorder struct {
	mock *MockCommandService_ExecuteClient
}

// NewMockCommandService_ExecuteClient creates a new mock instance.
func NewMockCommandService_ExecuteClient(ctrl *gomock.Controller) *MockCommandService_ExecuteClient {
	mock := &MockCommandService_ExecuteClient{ctrl: ctrl}
	mock.recorder = &MockCommandService_ExecuteClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandService_ExecuteClient) EXPECT() *MockCommandService_ExecuteClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockCommandService_ExecuteClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockCommandService_ExecuteClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockCommandService_ExecuteClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockCommandService_ExecuteClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockCommandService_ExecuteClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockCommandService_ExecuteClient)(nil).Context))
}

// Header mocks base method.
func (m *MockCommandService_ExecuteClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockCommandService_ExecuteClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockCommandService_ExecuteClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockCommandService_ExecuteClient) Recv() (*cloudproto1.CommandResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*cloudproto1.CommandResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockCommandService_ExecuteClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockCommandService_ExecuteClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockCommandService_ExecuteClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockCommandService_ExecuteClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockCommandService_ExecuteClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method.
func (m *MockCommandService_ExecuteClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockCommandService_ExecuteClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockCommandService_ExecuteClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockCommandService_ExecuteClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockCommandService_ExecuteClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockCommandService_ExecuteClient)(nil).Trailer))
}
