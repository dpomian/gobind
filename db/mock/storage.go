// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dpomian/gobind/db/sqlc (interfaces: Storage)
//
// Generated by this command:
//
//	mockgen -package mockdb -destination db/mock/storage.go github.com/dpomian/gobind/db/sqlc Storage
//

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/dpomian/gobind/db/sqlc"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// CreateNotebook mocks base method.
func (m *MockStorage) CreateNotebook(arg0 context.Context, arg1 db.CreateNotebookParams) (db.Notebook, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNotebook", arg0, arg1)
	ret0, _ := ret[0].(db.Notebook)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNotebook indicates an expected call of CreateNotebook.
func (mr *MockStorageMockRecorder) CreateNotebook(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNotebook", reflect.TypeOf((*MockStorage)(nil).CreateNotebook), arg0, arg1)
}

// DeleteNotebook mocks base method.
func (m *MockStorage) DeleteNotebook(arg0 context.Context, arg1 db.DeleteNotebookParams) (db.Notebook, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNotebook", arg0, arg1)
	ret0, _ := ret[0].(db.Notebook)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNotebook indicates an expected call of DeleteNotebook.
func (mr *MockStorageMockRecorder) DeleteNotebook(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNotebook", reflect.TypeOf((*MockStorage)(nil).DeleteNotebook), arg0, arg1)
}

// GetNotebook mocks base method.
func (m *MockStorage) GetNotebook(arg0 context.Context, arg1 uuid.UUID) (db.Notebook, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotebook", arg0, arg1)
	ret0, _ := ret[0].(db.Notebook)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotebook indicates an expected call of GetNotebook.
func (mr *MockStorageMockRecorder) GetNotebook(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotebook", reflect.TypeOf((*MockStorage)(nil).GetNotebook), arg0, arg1)
}

// ListNotebooks mocks base method.
func (m *MockStorage) ListNotebooks(arg0 context.Context, arg1 db.ListNotebooksParams) ([]db.Notebook, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNotebooks", arg0, arg1)
	ret0, _ := ret[0].([]db.Notebook)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNotebooks indicates an expected call of ListNotebooks.
func (mr *MockStorageMockRecorder) ListNotebooks(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNotebooks", reflect.TypeOf((*MockStorage)(nil).ListNotebooks), arg0, arg1)
}

// SearchNotebooks mocks base method.
func (m *MockStorage) SearchNotebooks(arg0 context.Context, arg1 string) ([]db.Notebook, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchNotebooks", arg0, arg1)
	ret0, _ := ret[0].([]db.Notebook)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchNotebooks indicates an expected call of SearchNotebooks.
func (mr *MockStorageMockRecorder) SearchNotebooks(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchNotebooks", reflect.TypeOf((*MockStorage)(nil).SearchNotebooks), arg0, arg1)
}

// UpdateNotebook mocks base method.
func (m *MockStorage) UpdateNotebook(arg0 context.Context, arg1 db.UpdateNotebookParams) (db.Notebook, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNotebook", arg0, arg1)
	ret0, _ := ret[0].(db.Notebook)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNotebook indicates an expected call of UpdateNotebook.
func (mr *MockStorageMockRecorder) UpdateNotebook(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNotebook", reflect.TypeOf((*MockStorage)(nil).UpdateNotebook), arg0, arg1)
}
