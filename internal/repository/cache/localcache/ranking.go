package localcache

import (
	"context"
	"errors"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/jellydator/ttlcache/v3"
)

type RankingLocalCache struct {
	cache.BaseRankingCache
	cache *ttlcache.Cache[string, []domain.Article]
}

// Get implements cache.RankingCache.
func (r *RankingLocalCache) Get(ctx context.Context) ([]domain.Article, error) {
	got := r.cache.Get(r.Key("article"))
	if got.IsExpired() {
		return []domain.Article{}, errors.New("local cache expired")
	}
	return got.Value(), nil
}

// Set implements cache.RankingCache.
func (r *RankingLocalCache) Set(ctx context.Context, arts []domain.Article) error {
	r.cache.Set(r.Key("article"), arts, ttlcache.DefaultTTL)
	return nil
}

func NewRankingLocalCache(
	cache *ttlcache.Cache[string, []domain.Article],
) *RankingLocalCache {
	return &RankingLocalCache{
		cache: cache,
	}
}
