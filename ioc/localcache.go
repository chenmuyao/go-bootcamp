package ioc

import (
	"time"

	intrCache "github.com/chenmuyao/go-bootcamp/interactive/repository/cache"
	intrLocalcache "github.com/chenmuyao/go-bootcamp/interactive/repository/cache/localcache"
	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache/localcache"
	"github.com/jellydator/ttlcache/v3"
)

func InitUserLocalCache() cache.UserCache {
	timeout := 15 * time.Minute
	cc := ttlcache.New(ttlcache.WithTTL[string, domain.User](timeout))
	go cc.Start()

	return localcache.NewUserLocalCache(cc)
}

func InitCodeLocalCache() cache.CodeCache {
	timeout := 10 * time.Minute
	code := ttlcache.New(
		ttlcache.WithTTL[string, string](timeout),
		ttlcache.WithDisableTouchOnHit[string, string](),
	)
	cnt := ttlcache.New(
		ttlcache.WithTTL[string, int](timeout),
		ttlcache.WithDisableTouchOnHit[string, int](),
	)
	go code.Start()
	go cnt.Start()

	return localcache.NewCodeLocalCache(code, cnt, timeout)
}

func InitTopArticlesCache() intrCache.TopArticlesCache {
	timeout := 15 * time.Second
	cc := ttlcache.New(
		ttlcache.WithTTL[string, []int64](timeout),
		ttlcache.WithDisableTouchOnHit[string, []int64](),
	)
	go cc.Start()

	return intrLocalcache.NewTopArticlesLocalCache(cc)
}

func InitRankingLocalCache() *localcache.RankingLocalCache {
	// NOTE: longer than redis, or never expires
	// If redis + sql down, can still use the expired local cache
	timeout := 5 * time.Minute
	cc := ttlcache.New(
		ttlcache.WithTTL[string, []domain.Article](timeout),
		ttlcache.WithDisableTouchOnHit[string, []domain.Article](),
	)
	go cc.Start()

	return localcache.NewRankingLocalCache(cc)
}
