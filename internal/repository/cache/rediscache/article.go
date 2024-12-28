package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

const (
	motifFirstPage                   = "first_page"
	motifContent                     = "content"
	articleFirstPageExpiryTime       = time.Minute
	articleContentPreCacheExpiryTime = 10 * time.Second
)

type ArticleRedisCache struct {
	cache.BaseArticleCache
	client redis.Cmdable
}

// Get implements cache.ArticleCache.
func (a *ArticleRedisCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	key := a.Key(motifContent, id)
	val, err := a.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var res domain.Article
	err = json.Unmarshal(val, &res)
	return res, err
}

// Set implements cache.ArticleCache.
func (a *ArticleRedisCache) Set(ctx context.Context, article domain.Article) error {
	key := a.Key(motifContent, article.ID)
	val, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return a.client.Set(ctx, key, val, articleContentPreCacheExpiryTime).Err()
}

// DelFirstPage implements cache.ArticleCache.
func (a *ArticleRedisCache) DelFirstPage(ctx context.Context, uid int64) error {
	key := a.Key(motifFirstPage, uid)
	return a.client.Del(ctx, key).Err()
}

// GetFirstPage implements cache.ArticleCache.
func (a *ArticleRedisCache) GetFirstPage(ctx context.Context, uid int64) ([]domain.Article, error) {
	key := a.Key(motifFirstPage, uid)
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
	key := a.Key(motifFirstPage, uid)
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
