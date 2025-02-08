package repository

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache/localcache"
)

//go:generate mockgen -source=./ranking.go -package=repomocks -destination=./mocks/ranking.mock.go
type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	rediscache cache.RankingCache
	localcache cache.RankingCache
}

// GetTopN implements RankingRepository.
func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	res, err := c.localcache.Get(ctx)
	if err == nil {
		return res, nil
	}
	res, err = c.rediscache.Get(ctx)
	if err != nil {
		return []domain.Article{}, nil
	}
	_ = c.localcache.Set(ctx, res)
	return res, nil
}

// ReplaceTopN implements RankingRepository.
func (c *CachedRankingRepository) ReplaceTopN(
	ctx context.Context,
	arts []domain.Article,
) error {
	_ = c.localcache.Set(ctx, arts)
	return c.rediscache.Set(ctx, arts)
}

func NewCachedRankingRepository(
	rediscache cache.RankingCache,
	localcache *localcache.RankingLocalCache,
) RankingRepository {
	return &CachedRankingRepository{
		rediscache: rediscache,
		localcache: localcache,
	}
}
