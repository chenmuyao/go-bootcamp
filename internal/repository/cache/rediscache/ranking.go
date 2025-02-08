package rediscache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/redis/go-redis/v9"
)

type RankingRedisCache struct {
	cache.BaseRankingCache
	client     redis.Cmdable
	expiration time.Duration
}

// Get implements cache.RankingCache.
func (r *RankingRedisCache) Get(ctx context.Context) ([]domain.Article, error) {
	res, err := r.client.Get(ctx, r.Key("article")).Bytes()
	if err != nil {
		return []domain.Article{}, err
	}

	var arts []domain.Article
	err = json.Unmarshal(res, &arts)
	if err != nil {
		return []domain.Article{}, err
	}
	return arts, nil
}

// Set implements cache.RankingCache.
func (r *RankingRedisCache) Set(ctx context.Context, arts []domain.Article) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.Key("article"), val, r.expiration).Err()
}

func NewRankingRedisCache(client redis.Cmdable) cache.RankingCache {
	return &RankingRedisCache{
		client:     client,
		expiration: 3 * time.Minute,
	}
}
