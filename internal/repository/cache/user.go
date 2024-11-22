package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/redis/go-redis/v9"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	// NOTE: Client or ClientCluster
	cmd redis.Cmdable
	// NOTE: only for user cache. Can be fixed here.
	expiration time.Duration

	// NOTE: keep expiration, key structure and marshal method locally
}

func NewUserCache(cmd redis.Cmdable) *UserCache {
	return &UserCache{
		cmd:        cmd,
		expiration: time.Minute * 15,
	}
}

func (c *UserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := c.Key(uid)
	// NOTE: Suppose using JSON to marshal
	data, err := c.cmd.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}

	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	return u, err
}

func (c *UserCache) Set(ctx context.Context, user domain.User) error {
	key := c.Key(user.ID)

	data, err := json.Marshal(&user)
	if err != nil {
		return err
	}

	return c.cmd.Set(ctx, key, &data, c.expiration).Err()
}

func (c *UserCache) Key(uid int64) string {
	return fmt.Sprintf("user:info:%d", uid)
}
