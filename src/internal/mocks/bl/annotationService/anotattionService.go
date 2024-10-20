// Code generated by MockGen. DO NOT EDIT.
// Source: bl/annotationService/anotattionService.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	models "annotater/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIAnotattionService is a mock of IAnotattionService interface.
type MockIAnotattionService struct {
	ctrl     *gomock.Controller
	recorder *MockIAnotattionServiceMockRecorder
}

// MockIAnotattionServiceMockRecorder is the mock recorder for MockIAnotattionService.
type MockIAnotattionServiceMockRecorder struct {
	mock *MockIAnotattionService
}

// NewMockIAnotattionService creates a new mock instance.
func NewMockIAnotattionService(ctrl *gomock.Controller) *MockIAnotattionService {
	mock := &MockIAnotattionService{ctrl: ctrl}
	mock.recorder = &MockIAnotattionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIAnotattionService) EXPECT() *MockIAnotattionServiceMockRecorder {
	return m.recorder
}

// AddAnottation mocks base method.
func (m *MockIAnotattionService) AddAnottation(anotattion *models.Markup) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAnottation", anotattion)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAnottation indicates an expected call of AddAnottation.
func (mr *MockIAnotattionServiceMockRecorder) AddAnottation(anotattion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAnottation", reflect.TypeOf((*MockIAnotattionService)(nil).AddAnottation), anotattion)
}

// DeleteAnotattion mocks base method.
func (m *MockIAnotattionService) DeleteAnotattion(id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAnotattion", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAnotattion indicates an expected call of DeleteAnotattion.
func (mr *MockIAnotattionServiceMockRecorder) DeleteAnotattion(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAnotattion", reflect.TypeOf((*MockIAnotattionService)(nil).DeleteAnotattion), id)
}

// GetAllAnottations mocks base method.
func (m *MockIAnotattionService) GetAllAnottations() ([]models.Markup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllAnottations")
	ret0, _ := ret[0].([]models.Markup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllAnottations indicates an expected call of GetAllAnottations.
func (mr *MockIAnotattionServiceMockRecorder) GetAllAnottations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllAnottations", reflect.TypeOf((*MockIAnotattionService)(nil).GetAllAnottations))
}

// GetAnottationByID mocks base method.
func (m *MockIAnotattionService) GetAnottationByID(id uint64) (*models.Markup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnottationByID", id)
	ret0, _ := ret[0].(*models.Markup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnottationByID indicates an expected call of GetAnottationByID.
func (mr *MockIAnotattionServiceMockRecorder) GetAnottationByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnottationByID", reflect.TypeOf((*MockIAnotattionService)(nil).GetAnottationByID), id)
}

// GetAnottationByUserID mocks base method.
func (m *MockIAnotattionService) GetAnottationByUserID(user_id uint64) ([]models.Markup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnottationByUserID", user_id)
	ret0, _ := ret[0].([]models.Markup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnottationByUserID indicates an expected call of GetAnottationByUserID.
func (mr *MockIAnotattionServiceMockRecorder) GetAnottationByUserID(user_id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnottationByUserID", reflect.TypeOf((*MockIAnotattionService)(nil).GetAnottationByUserID), user_id)
}
