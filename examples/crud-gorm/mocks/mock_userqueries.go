// Code generated by MockGen. DO NOT EDIT.
// Source: handlers/handler.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/go-fuego/fuego/examples/crud-gorm/models"
	gomock "github.com/golang/mock/gomock"
)

// MockUserQueryInterface is a mock of UserQueryInterface interface.
type MockUserQueryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserQueryInterfaceMockRecorder
}

// MockUserQueryInterfaceMockRecorder is the mock recorder for MockUserQueryInterface.
type MockUserQueryInterfaceMockRecorder struct {
	mock *MockUserQueryInterface
}

// NewMockUserQueryInterface creates a new mock instance.
func NewMockUserQueryInterface(ctrl *gomock.Controller) *MockUserQueryInterface {
	mock := &MockUserQueryInterface{ctrl: ctrl}
	mock.recorder = &MockUserQueryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserQueryInterface) EXPECT() *MockUserQueryInterfaceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserQueryInterface) CreateUser(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserQueryInterfaceMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserQueryInterface)(nil).CreateUser), user)
}

// DeleteUser mocks base method.
func (m *MockUserQueryInterface) DeleteUser(id uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockUserQueryInterfaceMockRecorder) DeleteUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserQueryInterface)(nil).DeleteUser), id)
}

// GetUserByEmail mocks base method.
func (m *MockUserQueryInterface) GetUserByEmail(email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserQueryInterfaceMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserQueryInterface)(nil).GetUserByEmail), email)
}

// GetUserByID mocks base method.
func (m *MockUserQueryInterface) GetUserByID(id uint) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserQueryInterfaceMockRecorder) GetUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserQueryInterface)(nil).GetUserByID), id)
}

// GetUsers mocks base method.
func (m *MockUserQueryInterface) GetUsers() ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsers")
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsers indicates an expected call of GetUsers.
func (mr *MockUserQueryInterfaceMockRecorder) GetUsers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsers", reflect.TypeOf((*MockUserQueryInterface)(nil).GetUsers))
}

// UpdateUser mocks base method.
func (m *MockUserQueryInterface) UpdateUser(user *models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockUserQueryInterfaceMockRecorder) UpdateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockUserQueryInterface)(nil).UpdateUser), user)
}