package repository

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

var (
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
)

// }}}
// {{{ Interface

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

// }}}
// {{{ Struct

type CachedCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cc cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: cc,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (c *CachedCodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CachedCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
