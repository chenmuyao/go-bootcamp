package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	cachemocks "github.com/chenmuyao/go-bootcamp/internal/repository/cache/mocks"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	daomocks "github.com/chenmuyao/go-bootcamp/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindByID(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx context.Context
		uid int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "find from cache ok",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				cache := cachemocks.NewMockUserCache(ctrl)
				cache.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					ID:       int64(123),
					Email:    "ok@test.com",
					Password: "pass",
					Phone:    "12345",
					Name:     "ok",
					Birthday: time.UnixMilli(123),
					Profile:  "my profile",
				}, nil)
				return nil, cache
			},
			ctx: context.Background(),
			uid: int64(123),
			wantUser: domain.User{
				ID:       int64(123),
				Email:    "ok@test.com",
				Password: "pass",
				Phone:    "12345",
				Name:     "ok",
				Birthday: time.UnixMilli(123),
				Profile:  "my profile",
			},
		},
		{
			name: "find from db ok",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				cc := cachemocks.NewMockUserCache(ctrl)
				cc.EXPECT().
					Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				cc.EXPECT().Set(gomock.Any(), domain.User{
					ID:       int64(123),
					Email:    "ok@test.com",
					Password: "pass",
					Phone:    "12345",
					Name:     "ok",
					Birthday: time.UnixMilli(123),
					Profile:  "my profile",
				}).Return(nil)
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindByID(gomock.Any(), int64(123)).Return(dao.User{
					ID:       int64(123),
					Email:    sql.NullString{String: "ok@test.com", Valid: true},
					Password: "pass",
					Phone:    sql.NullString{String: "12345", Valid: true},
					Ctime:    111,
					Utime:    222,
					Name:     "ok",
					Birthday: 123,
					Profile:  "my profile",
				}, nil)
				return d, cc
			},
			ctx: context.Background(),
			uid: int64(123),
			wantUser: domain.User{
				ID:       int64(123),
				Email:    "ok@test.com",
				Password: "pass",
				Phone:    "12345",
				Name:     "ok",
				Birthday: time.UnixMilli(123),
				Profile:  "my profile",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			dao, cache := tc.mock(ctrl)
			repo := NewUserRepository(dao, cache)

			resUser, err := repo.FindByID(tc.ctx, tc.uid)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantUser, resUser)
			if err := waitForMocks(context.Background(), ctrl); err != nil {
				t.Error(err)
			}
		})
	}
}

func waitForMocks(ctx context.Context, ctrl *gomock.Controller) error {
	ticker := time.NewTicker(1 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(3 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			if ctrl.Satisfied() {
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timeout waiting for mocks to be satisfied")
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		}
	}
}
