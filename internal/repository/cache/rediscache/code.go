package rediscache

import (
	"context"
	_ "embed"
	"errors"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode string
)

type CodeRedisCache struct {
	cache.BaseCodeCache
	cmd redis.Cmdable
}

func NewCodeRedisCache(cmd redis.Cmdable) cache.CodeCache {
	return &CodeRedisCache{
		cmd: cmd,
	}
}

func (c *CodeRedisCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		// Redis error
		return err
	}

	switch res {
	case -2:
		return errors.New("verification code exists but has no expiration date")
	case -1:
		return cache.ErrCodeSendTooMany
	default:
		return nil
	}
}

func (c *CodeRedisCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.Key(biz, phone)}, code).Int()
	if err != nil {
		// Redis error
		return false, err
	}

	switch res {
	case -2:
		return false, nil
	case -1:
		return false, cache.ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}
