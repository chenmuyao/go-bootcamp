package repository

import (
	"context"
	"testing"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	daomocks "github.com/chenmuyao/go-bootcamp/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCachedArticleRepository_Sync(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.ArticleReaderDAO, dao.ArticleAuthorDAO)

		article domain.Article
		wantID  int64
		wantErr error
	}{
		{
			name: "sync success",
			mock: func(ctrl *gomock.Controller) (dao.ArticleReaderDAO, dao.ArticleAuthorDAO) {
				adao := daomocks.NewMockArticleAuthorDAO(ctrl)
				adao.EXPECT().Create(gomock.Any(), dao.Article{
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
				}).Return(int64(1), nil)

				rdao := daomocks.NewMockArticleReaderDAO(ctrl)
				rdao.EXPECT().Upsert(gomock.Any(), dao.Article{
					ID:       1,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
				}).Return(nil)
				return rdao, adao
			},
			article: domain.Article{
				Title:   "my title",
				Content: "my content",
				Author: domain.Author{
					ID: 123,
				},
			},
			wantID: 1,
		},
		{
			name: "edit and sync success",
			mock: func(ctrl *gomock.Controller) (dao.ArticleReaderDAO, dao.ArticleAuthorDAO) {
				adao := daomocks.NewMockArticleAuthorDAO(ctrl)
				adao.EXPECT().UpdateByID(gomock.Any(), dao.Article{
					ID:       2,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
				}).Return(nil)

				rdao := daomocks.NewMockArticleReaderDAO(ctrl)
				rdao.EXPECT().Upsert(gomock.Any(), dao.Article{
					ID:       2,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
				}).Return(nil)
				return rdao, adao
			},
			article: domain.Article{
				ID:      2,
				Title:   "my title",
				Content: "my content",
				Author: domain.Author{
					ID: 123,
				},
			},
			wantID: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			readerDAO, authorDAO := tc.mock(ctrl)
			repo := NewArticleRepositoryV2(readerDAO, authorDAO)
			id, err := repo.SyncV1(context.TODO(), tc.article)
			assert.Equal(t, tc.wantID, id)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
