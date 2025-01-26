package localcache

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/jellydator/ttlcache/v3"
)

type TopArticlesLocalCache struct {
	cache *ttlcache.Cache[string, []domain.ArticleInteractive]
}

// GetTopLikedArticles implements cache.TopArticlesCache.
func (t *TopArticlesLocalCache) GetTopLikedArticles(
	ctx context.Context,
) ([]domain.ArticleInteractive, error) {
	res := t.cache.Get(t.key())
	return res.Value(), nil
}

// SetTopLikedArticles implements cache.TopArticlesCache.
func (t *TopArticlesLocalCache) SetTopLikedArticles(
	ctx context.Context,
	articles []domain.ArticleInteractive,
) error {
	_ = t.cache.Set(t.key(), articles, ttlcache.DefaultTTL)
	return nil
}

func (t *TopArticlesLocalCache) key() string {
	return "top_liked_articles"
}

func NewTopArticlesLocalCache(
	cache *ttlcache.Cache[string, []domain.ArticleInteractive],
) cache.TopArticlesCache {
	return &TopArticlesLocalCache{cache: cache}
}
