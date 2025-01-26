package ioc

import (
	"time"

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

func InitTopArticlesCache() cache.TopArticlesCache {
	timeout := 15 * time.Second
	cc := ttlcache.New(
		ttlcache.WithTTL[string, []domain.ArticleInteractive](timeout),
		ttlcache.WithDisableTouchOnHit[string, []domain.ArticleInteractive](),
	)
	go cc.Start()

	return localcache.NewTopArticlesLocalCache(cc)
}
