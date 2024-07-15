// Code generated by MockGen. DO NOT EDIT.
// Source: questionnaire_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/terui-ryota/offer-item/internal/domain/model"
	boil "github.com/volatiletech/sqlboiler/v4/boil"
)

// MockQuestionnaireRepository is a mock of QuestionnaireRepository interface.
type MockQuestionnaireRepository struct {
	ctrl     *gomock.Controller
	recorder *MockQuestionnaireRepositoryMockRecorder
}

// MockQuestionnaireRepositoryMockRecorder is the mock recorder for MockQuestionnaireRepository.
type MockQuestionnaireRepositoryMockRecorder struct {
	mock *MockQuestionnaireRepository
}

// NewMockQuestionnaireRepository creates a new mock instance.
func NewMockQuestionnaireRepository(ctrl *gomock.Controller) *MockQuestionnaireRepository {
	mock := &MockQuestionnaireRepository{ctrl: ctrl}
	mock.recorder = &MockQuestionnaireRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuestionnaireRepository) EXPECT() *MockQuestionnaireRepositoryMockRecorder {
	return m.recorder
}

// BulkGet mocks base method.
func (m *MockQuestionnaireRepository) BulkGet(ctx context.Context, exec boil.ContextExecutor, ids []model.OfferItemID) (map[model.OfferItemID]model.Questionnaire, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BulkGet", ctx, exec, ids)
	ret0, _ := ret[0].(map[model.OfferItemID]model.Questionnaire)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BulkGet indicates an expected call of BulkGet.
func (mr *MockQuestionnaireRepositoryMockRecorder) BulkGet(ctx, exec, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BulkGet", reflect.TypeOf((*MockQuestionnaireRepository)(nil).BulkGet), ctx, exec, ids)
}

// Delete mocks base method.
func (m *MockQuestionnaireRepository) Delete(ctx context.Context, tx *sql.Tx, id model.OfferItemID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, tx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockQuestionnaireRepositoryMockRecorder) Delete(ctx, tx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockQuestionnaireRepository)(nil).Delete), ctx, tx, id)
}

// Get mocks base method.
func (m *MockQuestionnaireRepository) Get(ctx context.Context, exec boil.ContextExecutor, id model.OfferItemID, withLock bool) (*model.Questionnaire, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, exec, id, withLock)
	ret0, _ := ret[0].(*model.Questionnaire)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockQuestionnaireRepositoryMockRecorder) Get(ctx, exec, id, withLock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockQuestionnaireRepository)(nil).Get), ctx, exec, id, withLock)
}

// Save mocks base method.
func (m *MockQuestionnaireRepository) Save(ctx context.Context, tx *sql.Tx, questionnaire model.Questionnaire) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, tx, questionnaire)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockQuestionnaireRepositoryMockRecorder) Save(ctx, tx, questionnaire interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockQuestionnaireRepository)(nil).Save), ctx, tx, questionnaire)
}
