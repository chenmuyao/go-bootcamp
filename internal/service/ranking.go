package service

import (
	"context"
	"log"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
)

type RankingService interface {
	TopN(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
	intrSvc   InteractiveService
	artSvc    ArticleService
	batchSize int
	scoreFunc func(likeCnt int64, utime time.Time) float64
	n         int
}

// TopN implements RankingService.
func (b *BatchRankingService) TopN(ctx context.Context) ([]domain.Article, error) {
	arts, err := b.topN(ctx)
	if err != nil {
		return []domain.Article{}, nil
	}
	// Save results to the cache
	log.Println(arts)
	return arts, err
}

func (b *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	return []domain.Article{}, nil
}

func NewBatchRankingService(intrSvc InteractiveService, artSvc ArticleService) RankingService {
	return &BatchRankingService{
		intrSvc: intrSvc,
		artSvc:  artSvc,
	}
}
