// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/chen-keinan/beacon/pkg/utils (interfaces: FolderMgr)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFolderMgr is a mock of FolderMgr interface.
type MockFolderMgr struct {
	ctrl     *gomock.Controller
	recorder *MockFolderMgrMockRecorder
}

// MockFolderMgrMockRecorder is the mock recorder for MockFolderMgr.
type MockFolderMgrMockRecorder struct {
	mock *MockFolderMgr
}

// NewMockFolderMgr creates a new mock instance.
func NewMockFolderMgr(ctrl *gomock.Controller) *MockFolderMgr {
	mock := &MockFolderMgr{ctrl: ctrl}
	mock.recorder = &MockFolderMgrMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFolderMgr) EXPECT() *MockFolderMgrMockRecorder {
	return m.recorder
}

// CreateFolder mocks base method.
func (m *MockFolderMgr) CreateFolder(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFolder", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFolder indicates an expected call of CreateFolder.
func (mr *MockFolderMgrMockRecorder) CreateFolder(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFolder", reflect.TypeOf((*MockFolderMgr)(nil).CreateFolder), arg0)
}

// GetHomeFolder mocks base method.
func (m *MockFolderMgr) GetHomeFolder() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHomeFolder")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHomeFolder indicates an expected call of GetHomeFolder.
func (mr *MockFolderMgrMockRecorder) GetHomeFolder() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHomeFolder", reflect.TypeOf((*MockFolderMgr)(nil).GetHomeFolder))
}
