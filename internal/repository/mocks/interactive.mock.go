// Code generated by MockGen. DO NOT EDIT.
// Source: ./interactive.go
//
// Generated by this command:
//
//	mockgen -source=./interactive.go -package=repomocks -destination=./mocks/interactive.mock.go
//

// Package repomocks is a generated GoMock package.
package repomocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/chenmuyao/go-bootcamp/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockInteractiveRepository is a mock of InteractiveRepository interface.
type MockInteractiveRepository struct {
	ctrl     *gomock.Controller
	recorder *MockInteractiveRepositoryMockRecorder
	isgomock struct{}
}

// MockInteractiveRepositoryMockRecorder is the mock recorder for MockInteractiveRepository.
type MockInteractiveRepositoryMockRecorder struct {
	mock *MockInteractiveRepository
}

// NewMockInteractiveRepository creates a new mock instance.
func NewMockInteractiveRepository(ctrl *gomock.Controller) *MockInteractiveRepository {
	mock := &MockInteractiveRepository{ctrl: ctrl}
	mock.recorder = &MockInteractiveRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInteractiveRepository) EXPECT() *MockInteractiveRepositoryMockRecorder {
	return m.recorder
}

// AddCollectionItem mocks base method.
func (m *MockInteractiveRepository) AddCollectionItem(ctx context.Context, biz string, id, cid, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCollectionItem", ctx, biz, id, cid, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCollectionItem indicates an expected call of AddCollectionItem.
func (mr *MockInteractiveRepositoryMockRecorder) AddCollectionItem(ctx, biz, id, cid, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCollectionItem", reflect.TypeOf((*MockInteractiveRepository)(nil).AddCollectionItem), ctx, biz, id, cid, uid)
}

// BatchIncrReadCnt mocks base method.
func (m *MockInteractiveRepository) BatchIncrReadCnt(ctx context.Context, bizs []string, bizIDs []int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchIncrReadCnt", ctx, bizs, bizIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchIncrReadCnt indicates an expected call of BatchIncrReadCnt.
func (mr *MockInteractiveRepositoryMockRecorder) BatchIncrReadCnt(ctx, bizs, bizIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchIncrReadCnt", reflect.TypeOf((*MockInteractiveRepository)(nil).BatchIncrReadCnt), ctx, bizs, bizIDs)
}

// BatchSetTopLike mocks base method.
func (m *MockInteractiveRepository) BatchSetTopLike(ctx context.Context, biz string, batchSize int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchSetTopLike", ctx, biz, batchSize)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchSetTopLike indicates an expected call of BatchSetTopLike.
func (mr *MockInteractiveRepositoryMockRecorder) BatchSetTopLike(ctx, biz, batchSize any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchSetTopLike", reflect.TypeOf((*MockInteractiveRepository)(nil).BatchSetTopLike), ctx, biz, batchSize)
}

// Collected mocks base method.
func (m *MockInteractiveRepository) Collected(ctx context.Context, biz string, bizID, uid int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Collected", ctx, biz, bizID, uid)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Collected indicates an expected call of Collected.
func (mr *MockInteractiveRepositoryMockRecorder) Collected(ctx, biz, bizID, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Collected", reflect.TypeOf((*MockInteractiveRepository)(nil).Collected), ctx, biz, bizID, uid)
}

// DecrLike mocks base method.
func (m *MockInteractiveRepository) DecrLike(ctx context.Context, biz string, id, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecrLike", ctx, biz, id, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DecrLike indicates an expected call of DecrLike.
func (mr *MockInteractiveRepositoryMockRecorder) DecrLike(ctx, biz, id, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecrLike", reflect.TypeOf((*MockInteractiveRepository)(nil).DecrLike), ctx, biz, id, uid)
}

// DeleteCollectionItem mocks base method.
func (m *MockInteractiveRepository) DeleteCollectionItem(ctx context.Context, biz string, id, cid, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCollectionItem", ctx, biz, id, cid, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCollectionItem indicates an expected call of DeleteCollectionItem.
func (mr *MockInteractiveRepositoryMockRecorder) DeleteCollectionItem(ctx, biz, id, cid, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCollectionItem", reflect.TypeOf((*MockInteractiveRepository)(nil).DeleteCollectionItem), ctx, biz, id, cid, uid)
}

// Get mocks base method.
func (m *MockInteractiveRepository) Get(ctx context.Context, biz string, bizID int64) (domain.Interactive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, biz, bizID)
	ret0, _ := ret[0].(domain.Interactive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockInteractiveRepositoryMockRecorder) Get(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockInteractiveRepository)(nil).Get), ctx, biz, bizID)
}

// GetByIDs mocks base method.
func (m *MockInteractiveRepository) GetByIDs(ctx context.Context, biz string, ids []int64) (map[int64]domain.Interactive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIDs", ctx, biz, ids)
	ret0, _ := ret[0].(map[int64]domain.Interactive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIDs indicates an expected call of GetByIDs.
func (mr *MockInteractiveRepositoryMockRecorder) GetByIDs(ctx, biz, ids any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDs", reflect.TypeOf((*MockInteractiveRepository)(nil).GetByIDs), ctx, biz, ids)
}

// GetTopLike mocks base method.
func (m *MockInteractiveRepository) GetTopLike(ctx context.Context, biz string, limit int) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopLike", ctx, biz, limit)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopLike indicates an expected call of GetTopLike.
func (mr *MockInteractiveRepositoryMockRecorder) GetTopLike(ctx, biz, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopLike", reflect.TypeOf((*MockInteractiveRepository)(nil).GetTopLike), ctx, biz, limit)
}

// IncrLike mocks base method.
func (m *MockInteractiveRepository) IncrLike(ctx context.Context, biz string, id, uid int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrLike", ctx, biz, id, uid)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrLike indicates an expected call of IncrLike.
func (mr *MockInteractiveRepositoryMockRecorder) IncrLike(ctx, biz, id, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrLike", reflect.TypeOf((*MockInteractiveRepository)(nil).IncrLike), ctx, biz, id, uid)
}

// IncrReadCnt mocks base method.
func (m *MockInteractiveRepository) IncrReadCnt(ctx context.Context, biz string, bizID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrReadCnt", ctx, biz, bizID)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrReadCnt indicates an expected call of IncrReadCnt.
func (mr *MockInteractiveRepositoryMockRecorder) IncrReadCnt(ctx, biz, bizID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrReadCnt", reflect.TypeOf((*MockInteractiveRepository)(nil).IncrReadCnt), ctx, biz, bizID)
}

// Liked mocks base method.
func (m *MockInteractiveRepository) Liked(ctx context.Context, biz string, bizID, uid int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Liked", ctx, biz, bizID, uid)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Liked indicates an expected call of Liked.
func (mr *MockInteractiveRepositoryMockRecorder) Liked(ctx, biz, bizID, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Liked", reflect.TypeOf((*MockInteractiveRepository)(nil).Liked), ctx, biz, bizID, uid)
}

// MustBatchGet mocks base method.
func (m *MockInteractiveRepository) MustBatchGet(ctx context.Context, biz string, bizIDs []int64) ([]domain.Interactive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MustBatchGet", ctx, biz, bizIDs)
	ret0, _ := ret[0].([]domain.Interactive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MustBatchGet indicates an expected call of MustBatchGet.
func (mr *MockInteractiveRepositoryMockRecorder) MustBatchGet(ctx, biz, bizIDs any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustBatchGet", reflect.TypeOf((*MockInteractiveRepository)(nil).MustBatchGet), ctx, biz, bizIDs)
}
