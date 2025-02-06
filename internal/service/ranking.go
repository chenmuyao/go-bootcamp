package service

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/chenmuyao/generique/gqueue"
	"github.com/chenmuyao/generique/gslice"
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
	offset := 0
	start := time.Now()

	// NOTE: Improve: DDL, only the last days
	ddl := start.Add(-7 * 24 * time.Hour)

	type Score struct {
		score float64
		art   domain.Article
	}
	topN := gqueue.NewPriorityQueue(b.n, func(src, dst Score) bool {
		// small heap
		return src.score < dst.score
	})

	for {
		arts, err := b.artSvc.ListPub(ctx, start, offset, b.batchSize)
		if err != nil {
			return []domain.Article{}, err
		}
		if len(arts) == 0 {
			break
		}
		ids := gslice.Map(arts, func(id int, src domain.Article) int64 {
			return src.ID
		})
		intrMap, err := b.intrSvc.GetByIDs(ctx, "article", ids)

		for _, art := range arts {
			intr := intrMap[art.ID]
			score := b.scoreFunc(intr.LikeCnt, art.Utime)
			topN.Enqueue(Score{
				score: score,
				art:   art,
			})
		}
		offset = offset + len(arts)
		// if len(arts) < b.batchSize {
		// 	// no more data
		// 	break
		// }

		// NOTE: improve with DDL
		if len(arts) < b.batchSize || arts[len(arts)-1].Utime.Before(ddl) {
			break
		}
	}
	res := make([]domain.Article, topN.Len())
	for i := topN.Len() - 1; i >= 0; i-- {
		el, err := topN.Dequeue()
		if err != nil {
			// log error
		}
		res[i] = el.art
	}
	return res, nil
}

func NewBatchRankingService(intrSvc InteractiveService, artSvc ArticleService) RankingService {
	return &BatchRankingService{
		intrSvc:   intrSvc,
		artSvc:    artSvc,
		batchSize: 100,
		n:         100,
		scoreFunc: func(likeCnt int64, utime time.Time) float64 {
			duration := time.Since(utime).Seconds()
			return float64(likeCnt-1) / math.Pow(duration+2, 1.5)
		},
	}
}
