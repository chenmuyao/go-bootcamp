package ioc

import (
	"github.com/chenmuyao/go-bootcamp/interactive/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: config.Cfg.Redis.Addr,
	})
}
