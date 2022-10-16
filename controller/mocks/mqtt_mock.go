// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/eclipse/paho.mqtt.golang (interfaces: Client)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// AddRoute mocks base method.
func (m *MockClient) AddRoute(arg0 string, arg1 mqtt.MessageHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddRoute", arg0, arg1)
}

// AddRoute indicates an expected call of AddRoute.
func (mr *MockClientMockRecorder) AddRoute(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRoute", reflect.TypeOf((*MockClient)(nil).AddRoute), arg0, arg1)
}

// Connect mocks base method.
func (m *MockClient) Connect() mqtt.Token {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(mqtt.Token)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockClientMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockClient)(nil).Connect))
}

// Disconnect mocks base method.
func (m *MockClient) Disconnect(arg0 uint) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Disconnect", arg0)
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockClientMockRecorder) Disconnect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockClient)(nil).Disconnect), arg0)
}

// IsConnected mocks base method.
func (m *MockClient) IsConnected() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsConnected")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsConnected indicates an expected call of IsConnected.
func (mr *MockClientMockRecorder) IsConnected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsConnected", reflect.TypeOf((*MockClient)(nil).IsConnected))
}

// IsConnectionOpen mocks base method.
func (m *MockClient) IsConnectionOpen() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsConnectionOpen")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsConnectionOpen indicates an expected call of IsConnectionOpen.
func (mr *MockClientMockRecorder) IsConnectionOpen() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsConnectionOpen", reflect.TypeOf((*MockClient)(nil).IsConnectionOpen))
}

// OptionsReader mocks base method.
func (m *MockClient) OptionsReader() mqtt.ClientOptionsReader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OptionsReader")
	ret0, _ := ret[0].(mqtt.ClientOptionsReader)
	return ret0
}

// OptionsReader indicates an expected call of OptionsReader.
func (mr *MockClientMockRecorder) OptionsReader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OptionsReader", reflect.TypeOf((*MockClient)(nil).OptionsReader))
}

// Publish mocks base method.
func (m *MockClient) Publish(arg0 string, arg1 byte, arg2 bool, arg3 interface{}) mqtt.Token {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(mqtt.Token)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockClientMockRecorder) Publish(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockClient)(nil).Publish), arg0, arg1, arg2, arg3)
}

// Subscribe mocks base method.
func (m *MockClient) Subscribe(arg0 string, arg1 byte, arg2 mqtt.MessageHandler) mqtt.Token {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1, arg2)
	ret0, _ := ret[0].(mqtt.Token)
	return ret0
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockClientMockRecorder) Subscribe(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockClient)(nil).Subscribe), arg0, arg1, arg2)
}

// SubscribeMultiple mocks base method.
func (m *MockClient) SubscribeMultiple(arg0 map[string]byte, arg1 mqtt.MessageHandler) mqtt.Token {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeMultiple", arg0, arg1)
	ret0, _ := ret[0].(mqtt.Token)
	return ret0
}

// SubscribeMultiple indicates an expected call of SubscribeMultiple.
func (mr *MockClientMockRecorder) SubscribeMultiple(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeMultiple", reflect.TypeOf((*MockClient)(nil).SubscribeMultiple), arg0, arg1)
}

// Unsubscribe mocks base method.
func (m *MockClient) Unsubscribe(arg0 ...string) mqtt.Token {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Unsubscribe", varargs...)
	ret0, _ := ret[0].(mqtt.Token)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockClientMockRecorder) Unsubscribe(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockClient)(nil).Unsubscribe), arg0...)
}
