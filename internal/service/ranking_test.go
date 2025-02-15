package service

import (
	"context"
	"testing"
	"time"

	intrDomain "github.com/chenmuyao/go-bootcamp/interactive/domain"
	"github.com/chenmuyao/go-bootcamp/interactive/service"
	intrsvcmocks "github.com/chenmuyao/go-bootcamp/interactive/service/mocks"
	"github.com/chenmuyao/go-bootcamp/internal/domain"
	svcmocks "github.com/chenmuyao/go-bootcamp/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBatchRankingService_TopN(t *testing.T) {
	now := time.Now()
	batchSize := 2
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (service.InteractiveService, ArticleService)

		wantArts []domain.Article
		wantErr  error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) (service.InteractiveService, ArticleService) {
				intrSvc := intrsvcmocks.NewMockInteractiveService(ctrl)
				artSvc := svcmocks.NewMockArticleService(ctrl)

				// 1st batch
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 0, 2).Return([]domain.Article{
					{ID: 1, Utime: now},
					{ID: 2, Utime: now},
				}, nil)
				intrSvc.EXPECT().
					GetByIDs(gomock.Any(), "article", []int64{1, 2}).
					Return(map[int64]intrDomain.Interactive{
						1: {LikeCnt: 1},
						2: {LikeCnt: 2},
					}, nil)

				// 2nd batch
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 2, 2).Return([]domain.Article{
					{ID: 3, Utime: now},
					{ID: 4, Utime: now},
				}, nil)
				intrSvc.EXPECT().
					GetByIDs(gomock.Any(), "article", []int64{3, 4}).
					Return(map[int64]intrDomain.Interactive{
						3: {LikeCnt: 3},
						4: {LikeCnt: 4},
					}, nil)

				// 3rd batch
				artSvc.EXPECT().
					ListPub(gomock.Any(), gomock.Any(), 4, 2).
					Return([]domain.Article{}, nil)

				return intrSvc, artSvc
			},
			wantErr: nil,
			wantArts: []domain.Article{
				{ID: 4, Utime: now},
				{ID: 3, Utime: now},
				{ID: 2, Utime: now},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			intrSvc, artSvc := tc.mock(ctrl)
			svc := &BatchRankingService{
				intrSvc:   intrSvc,
				artSvc:    artSvc,
				batchSize: batchSize,
				scoreFunc: func(likeCnt int64, utime time.Time) float64 {
					return float64(likeCnt)
				},
				n: 3,
			}
			arts, err := svc.topN(context.Background())
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantArts, arts)
		})
	}
}
