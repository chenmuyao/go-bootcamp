package rediscache

import (
	"context"
	"fmt"
	"testing"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	redismock "github.com/chenmuyao/go-bootcamp/internal/repository/cache/rediscache/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCodeCache_Set(t *testing.T) {
	keyFunc := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s:%s", biz, phone)
	}
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "set success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				res.EXPECT().
					Eval(gomock.Any(), luaSetCode, []string{keyFunc("test", "1234")}, "654321").
					Return(cmd)

				return res
			},
			ctx:     context.Background(),
			biz:     "test",
			phone:   "1234",
			code:    "654321",
			wantErr: nil,
		},
		{
			name: "redis error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(redis.ErrClosed)
				cmd.SetVal(int64(0))
				res.EXPECT().
					Eval(gomock.Any(), luaSetCode, []string{keyFunc("test", "1234")}, "654321").
					Return(cmd)

				return res
			},
			ctx:     context.Background(),
			biz:     "test",
			phone:   "1234",
			code:    "654321",
			wantErr: redis.ErrClosed,
		},
		{
			name: "no expiration date error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-2))
				res.EXPECT().
					Eval(gomock.Any(), luaSetCode, []string{keyFunc("test", "1234")}, "654321").
					Return(cmd)

				return res
			},
			ctx:     context.Background(),
			biz:     "test",
			phone:   "1234",
			code:    "654321",
			wantErr: ErrNoCodeExp,
		},
		{
			name: "too many request",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := redismock.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-1))
				res.EXPECT().
					Eval(gomock.Any(), luaSetCode, []string{keyFunc("test", "1234")}, "654321").
					Return(cmd)

				return res
			},
			ctx:     context.Background(),
			biz:     "test",
			phone:   "1234",
			code:    "654321",
			wantErr: cache.ErrCodeSendTooMany,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCodeRedisCache(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
