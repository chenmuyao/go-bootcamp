package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

// {{{ Consts

const defaultCacheExpiration = time.Minute * 15

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type UserRedisCache struct {
	cache.BaseUserCache
	// NOTE: Client or ClientCluster
	cmd redis.Cmdable
	// NOTE: only for user cache. Can be fixed here.
	expiration time.Duration

	// NOTE: keep expiration, key structure and marshal method locally
}

func NewUserRedisCache(cmd redis.Cmdable) cache.UserCache {
	return &UserRedisCache{
		cmd:        cmd,
		expiration: defaultCacheExpiration,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

// BatchGet implements cache.UserCache.
func (c *UserRedisCache) BatchGet(ctx context.Context, uids []int64) ([]domain.User, error) {
	res := make([]domain.User, 0, len(uids))
	for _, uid := range uids {
		u, err := c.Get(ctx, uid)
		if err != nil {
			return []domain.User{}, err
		}
		res = append(res, u)
	}
	return res, nil
}

// BatchSet implements cache.UserCache.
func (c *UserRedisCache) BatchSet(ctx context.Context, users []domain.User) error {
	var err error
	for _, user := range users {
		er := c.Set(ctx, user)
		if er != nil {
			// log the error
			err = er
		}
	}
	return err
}

func (c *UserRedisCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.Key(uid)
	// NOTE: Suppose using JSON to marshal
	data, err := c.cmd.Get(ctx, key).Result()
	if err == redis.Nil {
		return domain.User{}, cache.ErrKeyNotExist
	}
	if err != nil {
		return domain.User{}, err
	}

	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (c *UserRedisCache) Set(ctx context.Context, user domain.User) error {
	key := c.Key(user.ID)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.cmd.Set(ctx, key, data, c.expiration).Err()
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
