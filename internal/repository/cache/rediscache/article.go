package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

const articleFirstPageExpiryTime = time.Minute

type ArticleRedisCache struct {
	cache.BaseArticleCache
	client redis.Cmdable
}

// DelFirstPage implements cache.ArticleCache.
func (a *ArticleRedisCache) DelFirstPage(ctx context.Context, uid int64) error {
	key := a.Key(uid)
	return a.client.Del(ctx, key).Err()
}

// GetFirstPage implements cache.ArticleCache.
func (a *ArticleRedisCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	key := a.Key(uid)
	val, err := a.client.Get(ctx, key).Bytes()
	// val, err := a.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

// SetFirstPage implements cache.ArticleCache.
func (a *ArticleRedisCache) SetFirstPage(
	ctx context.Context,
	uid int64,
	articles []domain.Article,
) error {
	for i := range articles {
		articles[i].Content = articles[i].Abstract()
	}
	key := a.Key(uid)
	val, err := json.Marshal(articles)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, key, val, articleFirstPageExpiryTime).Err()
}

func NewArticleRedisCache(client redis.Cmdable) cache.ArticleCache {
	return &ArticleRedisCache{
		client: client,
	}
}
