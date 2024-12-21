package service

import (
	"context"
	"errors"
	"testing"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	repomocks "github.com/chenmuyao/go-bootcamp/internal/repository/mocks"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (repository.ArticleReaderRepository, repository.ArticleAuthorRepository)

		article domain.Article

		wantId  int64
		wantErr error
	}{
		{
			name: "publish success",
			mock: func(ctrl *gomock.Controller) (repository.ArticleReaderRepository, repository.ArticleAuthorRepository) {
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				}).Return(int64(1), nil)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      1,
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				})

				return readerRepo, authorRepo
			},
			article: domain.Article{
				Title:   "my title",
				Content: "my content",
				Author: domain.Author{
					ID: 123,
				},
			},
			wantId: 1,
		},
		{
			name: "edit publish success",
			mock: func(ctrl *gomock.Controller) (repository.ArticleReaderRepository, repository.ArticleAuthorRepository) {
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Update(gomock.Any(), domain.Article{
					ID:      2,
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				})
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      2,
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				})

				return readerRepo, authorRepo
			},
			article: domain.Article{
				ID:      2,
				Title:   "my title",
				Content: "my content",
				Author: domain.Author{
					ID: 123,
				},
				Status: domain.ArticleStatusPublished,
			},
			wantId: 2,
		},
		{
			name: "created but failed to publish, retry ok",
			mock: func(ctrl *gomock.Controller) (repository.ArticleReaderRepository, repository.ArticleAuthorRepository) {
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				}).Return(int64(1), nil)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      1,
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				}).Return(errors.New("publish error"))
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      1,
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				})

				return readerRepo, authorRepo
			},
			article: domain.Article{
				Title:   "my title",
				Content: "my content",
				Author: domain.Author{
					ID: 123,
				},
			},
			wantId:  1,
			wantErr: nil,
		},
		{
			name: "created but failed to publish, retry failed",
			mock: func(ctrl *gomock.Controller) (repository.ArticleReaderRepository, repository.ArticleAuthorRepository) {
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				}).Return(int64(1), nil)
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)
				readerRepo.EXPECT().Save(gomock.Any(), domain.Article{
					ID:      1,
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				}).Return(errors.New("publish error")).Times(publishMaxRetry)

				return readerRepo, authorRepo
			},
			article: domain.Article{
				Title:   "my title",
				Content: "my content",
				Author: domain.Author{
					ID: 123,
				},
			},
			wantId:  1,
			wantErr: ErrPublish,
		},
		{
			name: "failed to create",
			mock: func(ctrl *gomock.Controller) (repository.ArticleReaderRepository, repository.ArticleAuthorRepository) {
				authorRepo := repomocks.NewMockArticleAuthorRepository(ctrl)
				authorRepo.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID: 123,
					},
					Status: domain.ArticleStatusPublished,
				}).Return(int64(1), errors.New("create new article error"))
				readerRepo := repomocks.NewMockArticleReaderRepository(ctrl)

				return readerRepo, authorRepo
			},
			article: domain.Article{
				Title:   "my title",
				Content: "my content",
				Author: domain.Author{
					ID: 123,
				},
			},
			wantId:  0,
			wantErr: errors.New("create new article error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			readerRepo, authorRepo := tc.mock(ctrl)
			svc := NewArticleServiceV1(logger.NewZapLogger(zap.L()), readerRepo, authorRepo)
			id, err := svc.PublishV1(context.Background(), tc.article)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantId, id)
		})
	}
}
