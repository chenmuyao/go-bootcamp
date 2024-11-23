package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSentTooMany   = errors.New("too frequent send requests")
	ErrCodeVerifyTooMany = errors.New("too frequent verify requests")

	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

type CodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) *CodeCache {
	return &CodeCache{
		cmd: cmd,
	}
}

func (c *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		// Redis error
		return err
	}

	switch res {
	case -2:
		return errors.New("verification code exists but has no expiration date")
	case -1:
		return ErrCodeSentTooMany
	default:
		return nil
	}
}

func (c *CodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		// Redis error
		return false, err
	}

	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
