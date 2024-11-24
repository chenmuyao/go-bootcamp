package localcache

import (
	"context"
	"errors"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/jellydator/ttlcache/v3"
)

type UserLocalCache struct {
	cache.BaseUserCache
	cc *ttlcache.Cache[string, domain.User]
}

func NewUserLocalCache(cc *ttlcache.Cache[string, domain.User]) cache.UserCache {
	return &UserLocalCache{
		cc: cc,
	}
}

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
	res := c.cc.Set(key, user, ttlcache.DefaultTTL)
	if res == nil {
		return errors.New("cache set error")
	}
	return nil
}
