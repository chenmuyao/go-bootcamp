package localcache

import (
	"context"
	"errors"

	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/jellydator/ttlcache/v3"
)

type TopArticlesLocalCache struct {
	cache *ttlcache.Cache[string, []int64]
}

// GetTopLikedArticles implements cache.TopArticlesCache.
func (t *TopArticlesLocalCache) GetTopLikedArticles(
	ctx context.Context,
) ([]int64, error) {
	res := t.cache.Get(t.key())
	if res == nil {
		return nil, errors.New("No data")
	}
	return res.Value(), nil
}

// SetTopLikedArticles implements cache.TopArticlesCache.
func (t *TopArticlesLocalCache) SetTopLikedArticles(
	ctx context.Context,
	articles []int64,
) error {
	_ = t.cache.Set(t.key(), articles, ttlcache.DefaultTTL)
	return nil
}

func (t *TopArticlesLocalCache) key() string {
	return "top_liked_articles"
}

func NewTopArticlesLocalCache(
	cache *ttlcache.Cache[string, []int64],
) cache.TopArticlesCache {
	return &TopArticlesLocalCache{cache: cache}
}
