// Code generated by MockGen. DO NOT EDIT.
// Source: ./types.go
//
// Generated by this command:
//
//	mockgen -source=./types.go -package=cachemocks -destination=./mocks/cache.mock.go
//

// Package cachemocks is a generated GoMock package.
package cachemocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/chenmuyao/go-bootcamp/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockCodeCache is a mock of CodeCache interface.
type MockCodeCache struct {
	ctrl     *gomock.Controller
	recorder *MockCodeCacheMockRecorder
	isgomock struct{}
}

// MockCodeCacheMockRecorder is the mock recorder for MockCodeCache.
type MockCodeCacheMockRecorder struct {
	mock *MockCodeCache
}

// NewMockCodeCache creates a new mock instance.
func NewMockCodeCache(ctrl *gomock.Controller) *MockCodeCache {
	mock := &MockCodeCache{ctrl: ctrl}
	mock.recorder = &MockCodeCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCodeCache) EXPECT() *MockCodeCacheMockRecorder {
	return m.recorder
}

// Set mocks base method.
func (m *MockCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, biz, phone, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockCodeCacheMockRecorder) Set(ctx, biz, phone, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockCodeCache)(nil).Set), ctx, biz, phone, code)
}

// Verify mocks base method.
func (m *MockCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", ctx, biz, phone, code)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify indicates an expected call of Verify.
func (mr *MockCodeCacheMockRecorder) Verify(ctx, biz, phone, code any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockCodeCache)(nil).Verify), ctx, biz, phone, code)
}

// MockUserCache is a mock of UserCache interface.
type MockUserCache struct {
	ctrl     *gomock.Controller
	recorder *MockUserCacheMockRecorder
	isgomock struct{}
}

// MockUserCacheMockRecorder is the mock recorder for MockUserCache.
type MockUserCacheMockRecorder struct {
	mock *MockUserCache
}

// NewMockUserCache creates a new mock instance.
func NewMockUserCache(ctrl *gomock.Controller) *MockUserCache {
	mock := &MockUserCache{ctrl: ctrl}
	mock.recorder = &MockUserCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserCache) EXPECT() *MockUserCacheMockRecorder {
	return m.recorder
}

// BatchGet mocks base method.
func (m *MockUserCache) BatchGet(ctx context.Context, uids []int64) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchGet", ctx, uids)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchGet indicates an expected call of BatchGet.
func (mr *MockUserCacheMockRecorder) BatchGet(ctx, uids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchGet", reflect.TypeOf((*MockUserCache)(nil).BatchGet), ctx, uids)
}

// BatchSet mocks base method.
func (m *MockUserCache) BatchSet(ctx context.Context, users []domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSet", ctx, users)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchSet indicates an expected call of BatchSet.
func (mr *MockUserCacheMockRecorder) BatchSet(ctx, users any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSet", reflect.TypeOf((*MockUserCache)(nil).BatchSet), ctx, users)
}

// Get mocks base method.
func (m *MockUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, uid)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserCacheMockRecorder) Get(ctx, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserCache)(nil).Get), ctx, uid)
}

// Set mocks base method.
func (m *MockUserCache) Set(ctx context.Context, user domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockUserCacheMockRecorder) Set(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockUserCache)(nil).Set), ctx, user)
}

// MockArticleCache is a mock of ArticleCache interface.
type MockArticleCache struct {
	ctrl     *gomock.Controller
	recorder *MockArticleCacheMockRecorder
	isgomock struct{}
}

// MockArticleCacheMockRecorder is the mock recorder for MockArticleCache.
type MockArticleCacheMockRecorder struct {
	mock *MockArticleCache
}

// NewMockArticleCache creates a new mock instance.
func NewMockArticleCache(ctrl *gomock.Controller) *MockArticleCache {
	mock := &MockArticleCache{ctrl: ctrl}
	mock.recorder = &MockArticleCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleCache) EXPECT() *MockArticleCacheMockRecorder {
	return m.recorder
}

// BatchGetPub mocks base method.
func (m *MockArticleCache) BatchGetPub(ctx context.Context, ids []int64) ([]domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchGetPub", ctx, ids)
	ret0, _ := ret[0].([]domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchGetPub indicates an expected call of BatchGetPub.
func (mr *MockArticleCacheMockRecorder) BatchGetPub(ctx, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchGetPub", reflect.TypeOf((*MockArticleCache)(nil).BatchGetPub), ctx, ids)
}

// BatchSetPub mocks base method.
func (m *MockArticleCache) BatchSetPub(ctx context.Context, articles []domain.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSetPub", ctx, articles)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchSetPub indicates an expected call of BatchSetPub.
func (mr *MockArticleCacheMockRecorder) BatchSetPub(ctx, articles any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSetPub", reflect.TypeOf((*MockArticleCache)(nil).BatchSetPub), ctx, articles)
}

// DelFirstPage mocks base method.
func (m *MockArticleCache) DelFirstPage(ctx context.Context, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelFirstPage", ctx, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelFirstPage indicates an expected call of DelFirstPage.
func (mr *MockArticleCacheMockRecorder) DelFirstPage(ctx, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelFirstPage", reflect.TypeOf((*MockArticleCache)(nil).DelFirstPage), ctx, uid)
}

// Get mocks base method.
func (m *MockArticleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockArticleCacheMockRecorder) Get(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockArticleCache)(nil).Get), ctx, id)
}

// GetFirstPage mocks base method.
func (m *MockArticleCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFirstPage", ctx, uid)
	ret0, _ := ret[0].([]domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFirstPage indicates an expected call of GetFirstPage.
func (mr *MockArticleCacheMockRecorder) GetFirstPage(ctx, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFirstPage", reflect.TypeOf((*MockArticleCache)(nil).GetFirstPage), ctx, uid)
}

// GetPub mocks base method.
func (m *MockArticleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPub", ctx, id)
	ret0, _ := ret[0].(domain.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPub indicates an expected call of GetPub.
func (mr *MockArticleCacheMockRecorder) GetPub(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPub", reflect.TypeOf((*MockArticleCache)(nil).GetPub), ctx, id)
}

// Set mocks base method.
func (m *MockArticleCache) Set(ctx context.Context, article domain.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, article)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockArticleCacheMockRecorder) Set(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockArticleCache)(nil).Set), ctx, article)
}

// SetFirstPage mocks base method.
func (m *MockArticleCache) SetFirstPage(ctx context.Context, uid int64, articles []domain.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetFirstPage", ctx, uid, articles)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetFirstPage indicates an expected call of SetFirstPage.
func (mr *MockArticleCacheMockRecorder) SetFirstPage(ctx, uid, articles any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFirstPage", reflect.TypeOf((*MockArticleCache)(nil).SetFirstPage), ctx, uid, articles)
}

// SetPub mocks base method.
func (m *MockArticleCache) SetPub(ctx context.Context, article domain.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPub", ctx, article)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPub indicates an expected call of SetPub.
func (mr *MockArticleCacheMockRecorder) SetPub(ctx, article any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPub", reflect.TypeOf((*MockArticleCache)(nil).SetPub), ctx, article)
}

// MockInteractiveCache is a mock of InteractiveCache interface.
type MockInteractiveCache struct {
	ctrl     *gomock.Controller
	recorder *MockInteractiveCacheMockRecorder
	isgomock struct{}
}

// MockInteractiveCacheMockRecorder is the mock recorder for MockInteractiveCache.
type MockInteractiveCacheMockRecorder struct {
	mock *MockInteractiveCache
}

// NewMockInteractiveCache creates a new mock instance.
func NewMockInteractiveCache(ctrl *gomock.Controller) *MockInteractiveCache {
	mock := &MockInteractiveCache{ctrl: ctrl}
	mock.recorder = &MockInteractiveCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInteractiveCache) EXPECT() *MockInteractiveCacheMockRecorder {
	return m.recorder
}

// BatchSet mocks base method.
func (m *MockInteractiveCache) BatchSet(ctx context.Context, biz string, bizIDs []int64, intr []domain.Interactive) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSet", ctx, biz, bizIDs, intr)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchSet indicates an expected call of BatchSet.
func (mr *MockInteractiveCacheMockRecorder) BatchSet(ctx, biz, bizIDs, intr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSet", reflect.TypeOf((*MockInteractiveCache)(nil).BatchSet), ctx, biz, bizIDs, intr)
}

// DecrCollectCntIfPresent mocks base method.
func (m *MockInteractiveCache) DecrCollectCntIfPresent(ctx context.Context, biz string, bizID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecrCollectCntIfPresent", ctx, biz, bizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecrCollectCntIfPresent indicates an expected call of DecrCollectCntIfPresent.
func (mr *MockInteractiveCacheMockRecorder) DecrCollectCntIfPresent(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecrCollectCntIfPresent", reflect.TypeOf((*MockInteractiveCache)(nil).DecrCollectCntIfPresent), ctx, biz, bizID)
}

// DecrLikeCntIfPresent mocks base method.
func (m *MockInteractiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecrLikeCntIfPresent", ctx, biz, bizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecrLikeCntIfPresent indicates an expected call of DecrLikeCntIfPresent.
func (mr *MockInteractiveCacheMockRecorder) DecrLikeCntIfPresent(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecrLikeCntIfPresent", reflect.TypeOf((*MockInteractiveCache)(nil).DecrLikeCntIfPresent), ctx, biz, bizID)
}

// DecrLikeRank mocks base method.
func (m *MockInteractiveCache) DecrLikeRank(ctx context.Context, biz string, bizID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecrLikeRank", ctx, biz, bizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecrLikeRank indicates an expected call of DecrLikeRank.
func (mr *MockInteractiveCacheMockRecorder) DecrLikeRank(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecrLikeRank", reflect.TypeOf((*MockInteractiveCache)(nil).DecrLikeRank), ctx, biz, bizID)
}

// Get mocks base method.
func (m *MockInteractiveCache) Get(ctx context.Context, biz string, bizID int64) (domain.Interactive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, biz, bizID)
	ret0, _ := ret[0].(domain.Interactive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockInteractiveCacheMockRecorder) Get(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInteractiveCache)(nil).Get), ctx, biz, bizID)
}

// GetTopLikedIDs mocks base method.
func (m *MockInteractiveCache) GetTopLikedIDs(ctx context.Context, biz string, limit int64) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopLikedIDs", ctx, biz, limit)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopLikedIDs indicates an expected call of GetTopLikedIDs.
func (mr *MockInteractiveCacheMockRecorder) GetTopLikedIDs(ctx, biz, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopLikedIDs", reflect.TypeOf((*MockInteractiveCache)(nil).GetTopLikedIDs), ctx, biz, limit)
}

// IncrCollectCntIfPresent mocks base method.
func (m *MockInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrCollectCntIfPresent", ctx, biz, bizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrCollectCntIfPresent indicates an expected call of IncrCollectCntIfPresent.
func (mr *MockInteractiveCacheMockRecorder) IncrCollectCntIfPresent(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrCollectCntIfPresent", reflect.TypeOf((*MockInteractiveCache)(nil).IncrCollectCntIfPresent), ctx, biz, bizID)
}

// IncrLikeCntIfPresent mocks base method.
func (m *MockInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrLikeCntIfPresent", ctx, biz, bizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrLikeCntIfPresent indicates an expected call of IncrLikeCntIfPresent.
func (mr *MockInteractiveCacheMockRecorder) IncrLikeCntIfPresent(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrLikeCntIfPresent", reflect.TypeOf((*MockInteractiveCache)(nil).IncrLikeCntIfPresent), ctx, biz, bizID)
}

// IncrLikeRank mocks base method.
func (m *MockInteractiveCache) IncrLikeRank(ctx context.Context, biz string, bizID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrLikeRank", ctx, biz, bizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrLikeRank indicates an expected call of IncrLikeRank.
func (mr *MockInteractiveCacheMockRecorder) IncrLikeRank(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrLikeRank", reflect.TypeOf((*MockInteractiveCache)(nil).IncrLikeRank), ctx, biz, bizID)
}

// IncrReadCntIfPresent mocks base method.
func (m *MockInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrReadCntIfPresent", ctx, biz, bizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrReadCntIfPresent indicates an expected call of IncrReadCntIfPresent.
func (mr *MockInteractiveCacheMockRecorder) IncrReadCntIfPresent(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrReadCntIfPresent", reflect.TypeOf((*MockInteractiveCache)(nil).IncrReadCntIfPresent), ctx, biz, bizID)
}

// MustBatchGet mocks base method.
func (m *MockInteractiveCache) MustBatchGet(ctx context.Context, biz string, bizIDs []int64) ([]domain.Interactive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MustBatchGet", ctx, biz, bizIDs)
	ret0, _ := ret[0].([]domain.Interactive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MustBatchGet indicates an expected call of MustBatchGet.
func (mr *MockInteractiveCacheMockRecorder) MustBatchGet(ctx, biz, bizIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustBatchGet", reflect.TypeOf((*MockInteractiveCache)(nil).MustBatchGet), ctx, biz, bizIDs)
}

// Set mocks base method.
func (m *MockInteractiveCache) Set(ctx context.Context, biz string, bizID int64, intr domain.Interactive) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, biz, bizID, intr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockInteractiveCacheMockRecorder) Set(ctx, biz, bizID, intr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockInteractiveCache)(nil).Set), ctx, biz, bizID, intr)
}

// SetLikeToZSET mocks base method.
func (m *MockInteractiveCache) SetLikeToZSET(ctx context.Context, biz string, bizId, likeCnt int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLikeToZSET", ctx, biz, bizId, likeCnt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLikeToZSET indicates an expected call of SetLikeToZSET.
func (mr *MockInteractiveCacheMockRecorder) SetLikeToZSET(ctx, biz, bizId, likeCnt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLikeToZSET", reflect.TypeOf((*MockInteractiveCache)(nil).SetLikeToZSET), ctx, biz, bizId, likeCnt)
}

// MockTopArticlesCache is a mock of TopArticlesCache interface.
type MockTopArticlesCache struct {
	ctrl     *gomock.Controller
	recorder *MockTopArticlesCacheMockRecorder
	isgomock struct{}
}

// MockTopArticlesCacheMockRecorder is the mock recorder for MockTopArticlesCache.
type MockTopArticlesCacheMockRecorder struct {
	mock *MockTopArticlesCache
}

// NewMockTopArticlesCache creates a new mock instance.
func NewMockTopArticlesCache(ctrl *gomock.Controller) *MockTopArticlesCache {
	mock := &MockTopArticlesCache{ctrl: ctrl}
	mock.recorder = &MockTopArticlesCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopArticlesCache) EXPECT() *MockTopArticlesCacheMockRecorder {
	return m.recorder
}

// GetTopLikedArticles mocks base method.
func (m *MockTopArticlesCache) GetTopLikedArticles(ctx context.Context) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopLikedArticles", ctx)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopLikedArticles indicates an expected call of GetTopLikedArticles.
func (mr *MockTopArticlesCacheMockRecorder) GetTopLikedArticles(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopLikedArticles", reflect.TypeOf((*MockTopArticlesCache)(nil).GetTopLikedArticles), ctx)
}

// SetTopLikedArticles mocks base method.
func (m *MockTopArticlesCache) SetTopLikedArticles(ctx context.Context, articles []int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTopLikedArticles", ctx, articles)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTopLikedArticles indicates an expected call of SetTopLikedArticles.
func (mr *MockTopArticlesCacheMockRecorder) SetTopLikedArticles(ctx, articles any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTopLikedArticles", reflect.TypeOf((*MockTopArticlesCache)(nil).SetTopLikedArticles), ctx, articles)
}

// MockRankingCache is a mock of RankingCache interface.
type MockRankingCache struct {
	ctrl     *gomock.Controller
	recorder *MockRankingCacheMockRecorder
	isgomock struct{}
}

// MockRankingCacheMockRecorder is the mock recorder for MockRankingCache.
type MockRankingCacheMockRecorder struct {
	mock *MockRankingCache
}

// NewMockRankingCache creates a new mock instance.
func NewMockRankingCache(ctrl *gomock.Controller) *MockRankingCache {
	mock := &MockRankingCache{ctrl: ctrl}
	mock.recorder = &MockRankingCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRankingCache) EXPECT() *MockRankingCacheMockRecorder {
	return m.recorder
}

// Set mocks base method.
func (m *MockRankingCache) Set(ctx context.Context, arts []domain.Article) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, arts)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockRankingCacheMockRecorder) Set(ctx, arts any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRankingCache)(nil).Set), ctx, arts)
}
