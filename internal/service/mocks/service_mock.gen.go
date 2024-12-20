// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/service.go -destination=internal/service/mocks/service_mock.gen.go -package=svcmock
//

// Package svcmock is a generated GoMock package.
package svcmock

import (
	context "context"
	reflect "reflect"

	model "github.com/RIBorisov/GophKeeper/internal/model"
	storage "github.com/RIBorisov/GophKeeper/internal/storage"
	gomock "go.uber.org/mock/gomock"
)

// MockStoreI is a mock of StoreI interface.
type MockStoreI struct {
	ctrl     *gomock.Controller
	recorder *MockStoreIMockRecorder
}

// MockStoreIMockRecorder is the mock recorder for MockStoreI.
type MockStoreIMockRecorder struct {
	mock *MockStoreI
}

// NewMockStoreI creates a new mock instance.
func NewMockStoreI(ctrl *gomock.Controller) *MockStoreI {
	mock := &MockStoreI{ctrl: ctrl}
	mock.recorder = &MockStoreIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoreI) EXPECT() *MockStoreIMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockStoreI) Get(ctx context.Context, id string) (*storage.MetadataEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*storage.MetadataEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStoreIMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStoreI)(nil).Get), ctx, id)
}

// GetMany mocks base method.
func (m *MockStoreI) GetMany(ctx context.Context) ([]*storage.MetadataEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMany", ctx)
	ret0, _ := ret[0].([]*storage.MetadataEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMany indicates an expected call of GetMany.
func (mr *MockStoreIMockRecorder) GetMany(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMany", reflect.TypeOf((*MockStoreI)(nil).GetMany), ctx)
}

// GetUser mocks base method.
func (m *MockStoreI) GetUser(ctx context.Context, login string) (*storage.UserEntity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, login)
	ret0, _ := ret[0].(*storage.UserEntity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreIMockRecorder) GetUser(ctx, login any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStoreI)(nil).GetUser), ctx, login)
}

// Register mocks base method.
func (m *MockStoreI) Register(ctx context.Context, in model.UserCredentials) (model.UserID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, in)
	ret0, _ := ret[0].(model.UserID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockStoreIMockRecorder) Register(ctx, in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockStoreI)(nil).Register), ctx, in)
}

// Save mocks base method.
func (m *MockStoreI) Save(ctx context.Context, data model.Save) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, data)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockStoreIMockRecorder) Save(ctx, data any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockStoreI)(nil).Save), ctx, data)
}
