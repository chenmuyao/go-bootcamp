package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
)

var (
	ErrCodeSendTooMany   = errors.New("too frequent send requests")
	ErrCodeVerifyTooMany = errors.New("too frequent verify requests")
	ErrKeyNotExist       = errors.New("inexisted key")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
	// Key(biz, phone string) string
}

type BaseCodeCache struct{}

func (b *BaseCodeCache) Key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

type BaseUserCache struct{}

func (c *BaseUserCache) Key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}
