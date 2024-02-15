// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/deepakkamesh/medusa/controller/core (interfaces: MedusaCore)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	core "github.com/deepakkamesh/medusa/controller/core"
	gomock "github.com/golang/mock/gomock"
)

// MockMedusaCore is a mock of MedusaCore interface.
type MockMedusaCore struct {
	ctrl     *gomock.Controller
	recorder *MockMedusaCoreMockRecorder
}

// MockMedusaCoreMockRecorder is the mock recorder for MockMedusaCore.
type MockMedusaCoreMockRecorder struct {
	mock *MockMedusaCore
}

// NewMockMedusaCore creates a new mock instance.
func NewMockMedusaCore(ctrl *gomock.Controller) *MockMedusaCore {
	mock := &MockMedusaCore{ctrl: ctrl}
	mock.recorder = &MockMedusaCoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMedusaCore) EXPECT() *MockMedusaCoreMockRecorder {
	return m.recorder
}

// Action mocks base method.
func (m *MockMedusaCore) Action(arg0 []byte, arg1 byte, arg2 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Action", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Action indicates an expected call of Action.
func (mr *MockMedusaCoreMockRecorder) Action(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Action", reflect.TypeOf((*MockMedusaCore)(nil).Action), arg0, arg1, arg2)
}

// BoardConfig mocks base method.
func (m *MockMedusaCore) BoardConfig(arg0, arg1, arg2, arg3 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BoardConfig", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// BoardConfig indicates an expected call of BoardConfig.
func (mr *MockMedusaCoreMockRecorder) BoardConfig(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BoardConfig", reflect.TypeOf((*MockMedusaCore)(nil).BoardConfig), arg0, arg1, arg2, arg3)
}

// BuzzerOn mocks base method.
func (m *MockMedusaCore) BuzzerOn(arg0 []byte, arg1 bool, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuzzerOn", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// BuzzerOn indicates an expected call of BuzzerOn.
func (mr *MockMedusaCoreMockRecorder) BuzzerOn(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuzzerOn", reflect.TypeOf((*MockMedusaCore)(nil).BuzzerOn), arg0, arg1, arg2)
}

// CoreConfig mocks base method.
func (m *MockMedusaCore) CoreConfig() *core.Config {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CoreConfig")
	ret0, _ := ret[0].(*core.Config)
	return ret0
}

// CoreConfig indicates an expected call of CoreConfig.
func (mr *MockMedusaCoreMockRecorder) CoreConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CoreConfig", reflect.TypeOf((*MockMedusaCore)(nil).CoreConfig))
}

// Event mocks base method.
func (m *MockMedusaCore) Event() <-chan core.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Event")
	ret0, _ := ret[0].(<-chan core.Event)
	return ret0
}

// Event indicates an expected call of Event.
func (mr *MockMedusaCoreMockRecorder) Event() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Event", reflect.TypeOf((*MockMedusaCore)(nil).Event))
}

// GetBoardByAddr mocks base method.
func (m *MockMedusaCore) GetBoardByAddr(arg0 []byte) *core.Board {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardByAddr", arg0)
	ret0, _ := ret[0].(*core.Board)
	return ret0
}

// GetBoardByAddr indicates an expected call of GetBoardByAddr.
func (mr *MockMedusaCoreMockRecorder) GetBoardByAddr(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardByAddr", reflect.TypeOf((*MockMedusaCore)(nil).GetBoardByAddr), arg0)
}

// GetBoardByName mocks base method.
func (m *MockMedusaCore) GetBoardByName(arg0 string) *core.Board {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardByName", arg0)
	ret0, _ := ret[0].(*core.Board)
	return ret0
}

// GetBoardByName indicates an expected call of GetBoardByName.
func (mr *MockMedusaCoreMockRecorder) GetBoardByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardByName", reflect.TypeOf((*MockMedusaCore)(nil).GetBoardByName), arg0)
}

// GetBoardByRoom mocks base method.
func (m *MockMedusaCore) GetBoardByRoom(arg0 string) []core.Board {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoardByRoom", arg0)
	ret0, _ := ret[0].([]core.Board)
	return ret0
}

// GetBoardByRoom indicates an expected call of GetBoardByRoom.
func (mr *MockMedusaCoreMockRecorder) GetBoardByRoom(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoardByRoom", reflect.TypeOf((*MockMedusaCore)(nil).GetBoardByRoom), arg0)
}

// GetRelaybyPAddr mocks base method.
func (m *MockMedusaCore) GetRelaybyPAddr(arg0 []byte) *core.Relay {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRelaybyPAddr", arg0)
	ret0, _ := ret[0].(*core.Relay)
	return ret0
}

// GetRelaybyPAddr indicates an expected call of GetRelaybyPAddr.
func (mr *MockMedusaCoreMockRecorder) GetRelaybyPAddr(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRelaybyPAddr", reflect.TypeOf((*MockMedusaCore)(nil).GetRelaybyPAddr), arg0)
}

// LEDOn mocks base method.
func (m *MockMedusaCore) LEDOn(arg0 []byte, arg1 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LEDOn", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LEDOn indicates an expected call of LEDOn.
func (mr *MockMedusaCoreMockRecorder) LEDOn(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LEDOn", reflect.TypeOf((*MockMedusaCore)(nil).LEDOn), arg0, arg1)
}

// Light mocks base method.
func (m *MockMedusaCore) Light(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Light", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Light indicates an expected call of Light.
func (mr *MockMedusaCoreMockRecorder) Light(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Light", reflect.TypeOf((*MockMedusaCore)(nil).Light), arg0)
}

// RelayConfigMode mocks base method.
func (m *MockMedusaCore) RelayConfigMode(arg0 []byte, arg1 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RelayConfigMode", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RelayConfigMode indicates an expected call of RelayConfigMode.
func (mr *MockMedusaCoreMockRecorder) RelayConfigMode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RelayConfigMode", reflect.TypeOf((*MockMedusaCore)(nil).RelayConfigMode), arg0, arg1)
}

// Reset mocks base method.
func (m *MockMedusaCore) Reset(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Reset indicates an expected call of Reset.
func (mr *MockMedusaCoreMockRecorder) Reset(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockMedusaCore)(nil).Reset), arg0)
}

// StartCore mocks base method.
func (m *MockMedusaCore) StartCore() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartCore")
}

// StartCore indicates an expected call of StartCore.
func (mr *MockMedusaCoreMockRecorder) StartCore() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartCore", reflect.TypeOf((*MockMedusaCore)(nil).StartCore))
}

// Temp mocks base method.
func (m *MockMedusaCore) Temp(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Temp", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Temp indicates an expected call of Temp.
func (mr *MockMedusaCoreMockRecorder) Temp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Temp", reflect.TypeOf((*MockMedusaCore)(nil).Temp), arg0)
}

// Volt mocks base method.
func (m *MockMedusaCore) Volt(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Volt", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Volt indicates an expected call of Volt.
func (mr *MockMedusaCoreMockRecorder) Volt(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Volt", reflect.TypeOf((*MockMedusaCore)(nil).Volt), arg0)
}
