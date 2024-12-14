package localcache

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/jellydator/ttlcache/v3"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type UserLocalCache struct {
	cache.BaseUserCache
	cc *ttlcache.Cache[string, domain.User]
}

func NewUserLocalCache(cc *ttlcache.Cache[string, domain.User]) cache.UserCache {
	return &UserLocalCache{
		cc: cc,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (c *UserLocalCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.Key(uid)
	item := c.cc.Get(key)
	if item == nil {
		return domain.User{}, cache.ErrKeyNotExist
	}
	return item.Value(), nil
}

func (c *UserLocalCache) Set(ctx context.Context, user domain.User) error {
	key := c.Key(user.ID)
	c.cc.Set(key, user, ttlcache.DefaultTTL)
	return nil
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
