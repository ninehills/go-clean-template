// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/ninehills/go-webapp-template/internal/entity"
)

// MockUser is a mock of User interface.
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *MockUserMockRecorder
}

// MockUserMockRecorder is the mock recorder for MockUser.
type MockUserMockRecorder struct {
	mock *MockUser
}

// NewMockUser creates a new mock instance.
func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &MockUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUser) EXPECT() *MockUserMockRecorder {
	return m.recorder
}

// AuthenticationPassword mocks base method.
func (m *MockUser) AuthenticationPassword(ctx context.Context, username, password string) (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticationPassword", ctx, username, password)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AuthenticationPassword indicates an expected call of AuthenticationPassword.
func (mr *MockUserMockRecorder) AuthenticationPassword(ctx, username, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticationPassword", reflect.TypeOf((*MockUser)(nil).AuthenticationPassword), ctx, username, password)
}

// CacheGet mocks base method.
func (m *MockUser) CacheGet(ctx context.Context, username string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CacheGet", ctx, username)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CacheGet indicates an expected call of CacheGet.
func (mr *MockUserMockRecorder) CacheGet(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CacheGet", reflect.TypeOf((*MockUser)(nil).CacheGet), ctx, username)
}

// Create mocks base method.
func (m *MockUser) Create(ctx context.Context, in entity.User) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, in)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserMockRecorder) Create(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUser)(nil).Create), ctx, in)
}

// Delete mocks base method.
func (m *MockUser) Delete(ctx context.Context, username string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, username)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserMockRecorder) Delete(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUser)(nil).Delete), ctx, username)
}

// Get mocks base method.
func (m *MockUser) Get(ctx context.Context, username string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, username)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserMockRecorder) Get(ctx, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUser)(nil).Get), ctx, username)
}

// Query mocks base method.
func (m *MockUser) Query(ctx context.Context, p entity.PageQuery, o entity.OrderQuery, u entity.UserQuery) (entity.PageResult, []entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", ctx, p, o, u)
	ret0, _ := ret[0].(entity.PageResult)
	ret1, _ := ret[1].([]entity.User)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Query indicates an expected call of Query.
func (mr *MockUserMockRecorder) Query(ctx, p, o, u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockUser)(nil).Query), ctx, p, o, u)
}

// Update mocks base method.
func (m *MockUser) Update(ctx context.Context, in entity.User) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, in)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUserMockRecorder) Update(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUser)(nil).Update), ctx, in)
}
