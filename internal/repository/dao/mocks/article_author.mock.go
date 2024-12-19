// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/dao/article_author.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/dao/article_author.go -package=daomocks -destination=./internal/repository/dao/mocks/article_author.mock.go
//

// Package daomocks is a generated GoMock package.
package daomocks

import (
	context "context"
	reflect "reflect"

	dao "github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	gomock "go.uber.org/mock/gomock"
)

// MockArticleAuthorDAO is a mock of ArticleAuthorDAO interface.
type MockArticleAuthorDAO struct {
	ctrl     *gomock.Controller
	recorder *MockArticleAuthorDAOMockRecorder
	isgomock struct{}
}

// MockArticleAuthorDAOMockRecorder is the mock recorder for MockArticleAuthorDAO.
type MockArticleAuthorDAOMockRecorder struct {
	mock *MockArticleAuthorDAO
}

// NewMockArticleAuthorDAO creates a new mock instance.
func NewMockArticleAuthorDAO(ctrl *gomock.Controller) *MockArticleAuthorDAO {
	mock := &MockArticleAuthorDAO{ctrl: ctrl}
	mock.recorder = &MockArticleAuthorDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleAuthorDAO) EXPECT() *MockArticleAuthorDAOMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockArticleAuthorDAO) Create(ctx context.Context, article dao.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, article)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockArticleAuthorDAOMockRecorder) Create(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockArticleAuthorDAO)(nil).Create), ctx, article)
}

// UpdateByID mocks base method.
func (m *MockArticleAuthorDAO) UpdateByID(ctx context.Context, article dao.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateByID", ctx, article)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateByID indicates an expected call of UpdateByID.
func (mr *MockArticleAuthorDAOMockRecorder) UpdateByID(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateByID", reflect.TypeOf((*MockArticleAuthorDAO)(nil).UpdateByID), ctx, article)
}
