// Code generated by MockGen. DO NOT EDIT.
// Source: bl/documentService/documentService.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	models "annotater/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIDocumentService is a mock of IDocumentService interface.
type MockIDocumentService struct {
	ctrl     *gomock.Controller
	recorder *MockIDocumentServiceMockRecorder
}

// MockIDocumentServiceMockRecorder is the mock recorder for MockIDocumentService.
type MockIDocumentServiceMockRecorder struct {
	mock *MockIDocumentService
}

// NewMockIDocumentService creates a new mock instance.
func NewMockIDocumentService(ctrl *gomock.Controller) *MockIDocumentService {
	mock := &MockIDocumentService{ctrl: ctrl}
	mock.recorder = &MockIDocumentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDocumentService) EXPECT() *MockIDocumentServiceMockRecorder {
	return m.recorder
}

// CheckDocument mocks base method.
func (m *MockIDocumentService) CheckDocument(document models.DocumentMetaData) ([]models.Markup, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckDocument", document)
	ret0, _ := ret[0].([]models.Markup)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckDocument indicates an expected call of CheckDocument.
func (mr *MockIDocumentServiceMockRecorder) CheckDocument(document interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckDocument", reflect.TypeOf((*MockIDocumentService)(nil).CheckDocument), document)
}

// LoadDocument mocks base method.
func (m *MockIDocumentService) LoadDocument(document models.DocumentMetaData) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadDocument", document)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadDocument indicates an expected call of LoadDocument.
func (mr *MockIDocumentServiceMockRecorder) LoadDocument(document interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadDocument", reflect.TypeOf((*MockIDocumentService)(nil).LoadDocument), document)
}