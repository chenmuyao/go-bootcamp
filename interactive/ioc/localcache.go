package ioc

import (
	"time"

	intrCache "github.com/chenmuyao/go-bootcamp/interactive/repository/cache"
	intrLocalcache "github.com/chenmuyao/go-bootcamp/interactive/repository/cache/localcache"
	"github.com/jellydator/ttlcache/v3"
)

func InitTopArticlesCache() intrCache.TopArticlesCache {
	timeout := 15 * time.Second
	cc := ttlcache.New(
		ttlcache.WithTTL[string, []int64](timeout),
		ttlcache.WithDisableTouchOnHit[string, []int64](),
	)
	go cc.Start()

	return intrLocalcache.NewTopArticlesLocalCache(cc)
}
