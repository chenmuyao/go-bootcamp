// Code generated by MockGen. DO NOT EDIT.
// Source: ./article.go
//
// Generated by this command:
//
//	mockgen -source=./article.go -package=daomocks -destination=./mocks/article.mock.go
//

// Package daomocks is a generated GoMock package.
package daomocks

import (
	context "context"
	reflect "reflect"

	dao "github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	gomock "go.uber.org/mock/gomock"
)

// MockArticleDAO is a mock of ArticleDAO interface.
type MockArticleDAO struct {
	ctrl     *gomock.Controller
	recorder *MockArticleDAOMockRecorder
	isgomock struct{}
}

// MockArticleDAOMockRecorder is the mock recorder for MockArticleDAO.
type MockArticleDAOMockRecorder struct {
	mock *MockArticleDAO
}

// NewMockArticleDAO creates a new mock instance.
func NewMockArticleDAO(ctrl *gomock.Controller) *MockArticleDAO {
	mock := &MockArticleDAO{ctrl: ctrl}
	mock.recorder = &MockArticleDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleDAO) EXPECT() *MockArticleDAOMockRecorder {
	return m.recorder
}

// BatchGetPubByIDs mocks base method.
func (m *MockArticleDAO) BatchGetPubByIDs(ctx context.Context, ids []int64) ([]dao.PublishedArticle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchGetPubByIDs", ctx, ids)
	ret0, _ := ret[0].([]dao.PublishedArticle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchGetPubByIDs indicates an expected call of BatchGetPubByIDs.
func (mr *MockArticleDAOMockRecorder) BatchGetPubByIDs(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchGetPubByIDs", reflect.TypeOf((*MockArticleDAO)(nil).BatchGetPubByIDs), ctx, ids)
}

// GetByAuthor mocks base method.
func (m *MockArticleDAO) GetByAuthor(ctx context.Context, uid int64, offset, limit int) ([]dao.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAuthor", ctx, uid, offset, limit)
	ret0, _ := ret[0].([]dao.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAuthor indicates an expected call of GetByAuthor.
func (mr *MockArticleDAOMockRecorder) GetByAuthor(ctx, uid, offset, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAuthor", reflect.TypeOf((*MockArticleDAO)(nil).GetByAuthor), ctx, uid, offset, limit)
}

// GetByID mocks base method.
func (m *MockArticleDAO) GetByID(ctx context.Context, id int64) (dao.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(dao.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockArticleDAOMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockArticleDAO)(nil).GetByID), ctx, id)
}

// GetPubByID mocks base method.
func (m *MockArticleDAO) GetPubByID(ctx context.Context, id int64) (dao.PublishedArticle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPubByID", ctx, id)
	ret0, _ := ret[0].(dao.PublishedArticle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPubByID indicates an expected call of GetPubByID.
func (mr *MockArticleDAOMockRecorder) GetPubByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPubByID", reflect.TypeOf((*MockArticleDAO)(nil).GetPubByID), ctx, id)
}

// Insert mocks base method.
func (m *MockArticleDAO) Insert(ctx context.Context, article dao.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, article)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockArticleDAOMockRecorder) Insert(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockArticleDAO)(nil).Insert), ctx, article)
}

// Sync mocks base method.
func (m *MockArticleDAO) Sync(ctx context.Context, article dao.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sync", ctx, article)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync.
func (mr *MockArticleDAOMockRecorder) Sync(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockArticleDAO)(nil).Sync), ctx, article)
}

// Transaction mocks base method.
func (m *MockArticleDAO) Transaction(ctx context.Context, fn func(context.Context, any) (any, error)) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transaction", ctx, fn)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Transaction indicates an expected call of Transaction.
func (mr *MockArticleDAOMockRecorder) Transaction(ctx, fn any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction", reflect.TypeOf((*MockArticleDAO)(nil).Transaction), ctx, fn)
}

// UpdateByID mocks base method.
func (m *MockArticleDAO) UpdateByID(ctx context.Context, article dao.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateByID", ctx, article)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateByID indicates an expected call of UpdateByID.
func (mr *MockArticleDAOMockRecorder) UpdateByID(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateByID", reflect.TypeOf((*MockArticleDAO)(nil).UpdateByID), ctx, article)
}

// UpdateStatusByID mocks base method.
func (m *MockArticleDAO) UpdateStatusByID(ctx context.Context, model any, userID, articleID int64, status uint8) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatusByID", ctx, model, userID, articleID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStatusByID indicates an expected call of UpdateStatusByID.
func (mr *MockArticleDAOMockRecorder) UpdateStatusByID(ctx, model, userID, articleID, status any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatusByID", reflect.TypeOf((*MockArticleDAO)(nil).UpdateStatusByID), ctx, model, userID, articleID, status)
}

// Upsert mocks base method.
func (m *MockArticleDAO) Upsert(ctx context.Context, article dao.PublishedArticle) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upsert", ctx, article)
	ret0, _ := ret[0].(error)
	return ret0
}

// Upsert indicates an expected call of Upsert.
func (mr *MockArticleDAOMockRecorder) Upsert(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upsert", reflect.TypeOf((*MockArticleDAO)(nil).Upsert), ctx, article)
}
