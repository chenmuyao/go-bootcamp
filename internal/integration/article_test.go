package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/integration/startup"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Result[T any] struct {
	Data T      `json:"data"`
	Msg  string `json:"message"`
	Code int    `json:"code"`
}

type ArticleHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	rdb    redis.Cmdable
	server *gin.Engine
}

func (s *ArticleHandlerSuite) SetupSuite() {
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UserClaims{
			UID: 123,
		})
	})
	s.db = startup.InitDB()
	hdl := startup.InitArticleHandler(dao.NewArticleDAO(s.db))
	hdl.RegisterRoutes(s.server)
	s.rdb = startup.InitRedis()
}

func (s *ArticleHandlerSuite) TearDownTest() {
	log.Println("TRUNCATE")
	s.db.Exec("TRUNCATE articles")
	s.db.Exec("TRUNCATE published_articles")
	s.db.Exec("TRUNCATE users")
	s.db.Exec("TRUNCATE interactives")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.rdb.Del(ctx, "user:info:123")
}

func (s *ArticleHandlerSuite) TestEdit() {
	t := s.T()

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		// json article from frontend
		article web.ArticleEditReq

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name:   "Create a new post",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				// check that the article is saved into the DB
				var article dao.Article
				err := s.db.Where("author_id=?", 123).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Ctime > 789)
				article.Ctime = 0
				assert.True(t, article.Utime > 789)
				article.Utime = 0
				assert.Equal(t, dao.Article{
					ID:       1,
					Title:    "my super title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusUnpublished,
				}, article)
			},
			article: web.ArticleEditReq{
				Title:   "my super title",
				Content: "my content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: ginx.CodeOK,
				Data: 1,
			},
		},
		{
			name: "Edit a post",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					ID:       2,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the article is saved into the DB
				var article dao.Article
				err := s.db.Where("id = ?", 2).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Utime > 789)
				article.Utime = 0
				assert.Equal(t, dao.Article{
					ID:       2,
					Title:    "new title",
					Content:  "new content",
					AuthorID: 123,
					Status:   domain.ArticleStatusUnpublished,
					Ctime:    456,
				}, article)
			},
			article: web.ArticleEditReq{
				ID:      2,
				Title:   "new title",
				Content: "new content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: ginx.CodeOK,
				Data: 2,
			},
		},
		{
			name: "Edit a post of someone else",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					ID:       3,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the article is saved into the DB
				var article dao.Article
				err := s.db.Where("id = ?", 3).First(&article).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					ID:       3,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}, article)
			},
			article: web.ArticleEditReq{
				ID:      3,
				Title:   "new title",
				Content: "new content",
			},
			wantCode: http.StatusBadRequest,
			wantRes: Result[int64]{
				Code: ginx.CodeUserSide,
				Msg:  "article not found",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(tc.article)
			assert.NoError(t, err)
			req, err := http.NewRequest(
				http.MethodPost,
				"/articles/edit",
				bytes.NewReader(reqBody),
			)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			rec := httptest.NewRecorder()

			s.server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			var res Result[int64]
			err = json.NewDecoder(rec.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func (s *ArticleHandlerSuite) TestPublish() {
	t := s.T()

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		// json article from frontend
		article web.ArticlePublishReq

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name:   "Publish a new post",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				// check that the article is saved into the DB
				var article dao.Article
				err := s.db.Where("author_id=?", 123).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Ctime > 789)
				article.Ctime = 0
				assert.True(t, article.Utime > 789)
				article.Utime = 0
				assert.Equal(t, dao.Article{
					ID:       1,
					Title:    "my super title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
				}, article)

				var articlePub dao.PublishedArticle
				err = s.db.Where("author_id=?", 123).First(&articlePub).Error
				assert.NoError(t, err)
				assert.True(t, articlePub.Ctime > 789)
				articlePub.Ctime = 0
				assert.True(t, articlePub.Utime > 789)
				articlePub.Utime = 0
				assert.Equal(t, dao.PublishedArticle{
					ID:       1,
					Title:    "my super title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
				}, articlePub)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := s.rdb.Get(ctx, "article:published_content:1").Bytes()
				assert.NoError(t, err)
				var res domain.Article
				err = json.Unmarshal(val, &res)
				assert.NoError(t, err)
				assert.True(t, res.Ctime.UnixMilli() > 789)
				res.Ctime = time.Time{}
				assert.True(t, res.Utime.UnixMilli() > 789)
				res.Utime = time.Time{}
				assert.Equal(t, domain.Article{
					ID:      1,
					Title:   "my super title",
					Content: "my content",
					Author:  domain.Author{}, // Author is not set in the DB
					Status:  domain.ArticleStatusPublished,
				}, res)
				assert.NoError(t, s.rdb.Del(ctx, "article:published_content:1").Err())
			},
			article: web.ArticlePublishReq{
				Title:   "my super title",
				Content: "my content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: ginx.CodeOK,
				Data: 1,
			},
		},
		{
			name: "Edit and publish a post",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					ID:       22,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
				err = s.db.Create(dao.PublishedArticle{
					ID:       22,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the article is saved into the DB
				var article dao.Article
				err := s.db.Where("id = ?", 22).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Utime > 789)
				article.Utime = 0
				assert.Equal(t, dao.Article{
					ID:       22,
					Title:    "new title",
					Content:  "new content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
				}, article)

				var pubArticle dao.PublishedArticle
				err = s.db.Where("id = ?", 22).First(&pubArticle).Error
				assert.NoError(t, err)
				assert.True(t, pubArticle.Utime > 789)
				pubArticle.Utime = 0
				assert.Equal(t, dao.PublishedArticle{
					ID:       22,
					Title:    "new title",
					Content:  "new content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
				}, pubArticle)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := s.rdb.Get(ctx, "article:published_content:22").Bytes()
				assert.NoError(t, err)
				var res domain.Article
				err = json.Unmarshal(val, &res)
				assert.NoError(t, err)
				assert.True(t, res.Ctime.UnixMilli() > 456)
				res.Ctime = time.Time{}
				assert.True(t, res.Utime.UnixMilli() > 789)
				res.Utime = time.Time{}
				assert.Equal(t, domain.Article{
					ID:      22,
					Title:   "new title",
					Content: "new content",
					Author:  domain.Author{}, // Author is not set in the DB
					Status:  domain.ArticleStatusPublished,
				}, res)
				assert.NoError(t, s.rdb.Del(ctx, "article:published_content:22").Err())
			},
			article: web.ArticlePublishReq{
				ID:      22,
				Title:   "new title",
				Content: "new content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: ginx.CodeOK,
				Data: 22,
			},
		},
		{
			name: "Publish a post of someone else",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					ID:       23,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
				err = s.db.Create(dao.PublishedArticle{
					ID:       23,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the article is saved into the DB
				var article dao.Article
				err := s.db.Where("id = ?", 23).First(&article).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					ID:       23,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}, article)

				var pubArticle dao.Article
				err = s.db.Where("id = ?", 23).First(&pubArticle).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					ID:       23,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}, pubArticle)
			},
			article: web.ArticlePublishReq{
				ID:      23,
				Title:   "new title",
				Content: "new content",
			},
			wantCode: http.StatusBadRequest,
			wantRes: Result[int64]{
				Code: ginx.CodeUserSide,
				Msg:  "article not found",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(tc.article)
			assert.NoError(t, err)
			req, err := http.NewRequest(
				http.MethodPost,
				"/articles/publish",
				bytes.NewReader(reqBody),
			)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			rec := httptest.NewRecorder()

			s.server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			var res Result[int64]
			err = json.NewDecoder(rec.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func (s *ArticleHandlerSuite) TestWithdraw() {
	t := s.T()

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		// json article from frontend
		article web.ArticleWithdrawReq

		wantCode int
		wantRes  Result[int64]
	}{
		{
			name: "Withdraw a published post",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					ID:       31,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
				err = s.db.Create(dao.PublishedArticle{
					ID:       31,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the article is saved into the DB
				var article dao.Article
				err := s.db.Where("id = ?", 31).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, domain.ArticleStatusPrivate == article.Status)
				var pub dao.PublishedArticle
				err = s.db.Where("id = ?", 31).First(&pub).Error
				assert.NoError(t, err)
			},
			article: web.ArticleWithdrawReq{
				ID: 31,
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: ginx.CodeOK,
			},
		},
		{
			name: "Edit a post of someone else",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					ID:       32,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
				err = s.db.Create(dao.PublishedArticle{
					ID:       32,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the article is saved into the DB
				var article dao.Article
				err := s.db.Where("id = ?", 32).First(&article).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					ID:       32,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}, article)
			},
			article: web.ArticleWithdrawReq{
				ID: 32,
			},
			wantCode: http.StatusBadRequest,
			wantRes: Result[int64]{
				Code: ginx.CodeUserSide,
				Msg:  "article not found",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(tc.article)
			assert.NoError(t, err)
			req, err := http.NewRequest(
				http.MethodPost,
				"/articles/withdraw",
				bytes.NewReader(reqBody),
			)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			rec := httptest.NewRecorder()

			s.server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			var res Result[int64]
			err = json.NewDecoder(rec.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func (s *ArticleHandlerSuite) TestDetail() {
	t := s.T()

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		param string

		wantCode int
		wantRes  Result[web.ArticleVO]
	}{
		{
			name: "authors see detail of an article from db",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					ID:       41,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusUnpublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the cache is set
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := s.rdb.Get(ctx, "article:content:41").Bytes()
				assert.NoError(t, err)
				var res domain.Article
				err = json.Unmarshal(val, &res)
				assert.NoError(t, err)
				assert.Equal(t, domain.Article{
					ID:      41,
					Title:   "my title",
					Content: "my content",
					Author:  domain.Author{ID: 123},
					Status:  domain.ArticleStatusUnpublished,
					Ctime:   time.UnixMilli(456),
					Utime:   time.UnixMilli(789),
				}, res)
				assert.NoError(t, s.rdb.Del(ctx, "article:content:41").Err())
			},
			param:    "41",
			wantCode: http.StatusOK,
			wantRes: Result[web.ArticleVO]{
				Code: ginx.CodeOK,
				Data: web.ArticleVO{
					ID:      41,
					Title:   "my title",
					Content: "my content",
					Status:  domain.ArticleStatusUnpublished,
					Ctime:   time.UnixMilli(456).Format(time.DateTime),
					Utime:   time.UnixMilli(789).Format(time.DateTime),
				},
			},
		},
		{
			name: "authors see detail of an article from cache",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := json.Marshal(domain.Article{
					ID:      42,
					Title:   "my title",
					Content: "my content",
					Author:  domain.Author{ID: 123},
					Status:  domain.ArticleStatusUnpublished,
					Ctime:   time.UnixMilli(456),
					Utime:   time.UnixMilli(789),
				})
				assert.NoError(t, err)
				err = s.rdb.Set(ctx, "article:content:42", val, time.Second).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				assert.NoError(t, s.rdb.Del(ctx, "article:content:42").Err())
			},
			param:    "42",
			wantCode: http.StatusOK,
			wantRes: Result[web.ArticleVO]{
				Code: ginx.CodeOK,
				Data: web.ArticleVO{
					ID:      42,
					Title:   "my title",
					Content: "my content",
					Status:  domain.ArticleStatusUnpublished,
					Ctime:   time.UnixMilli(456).Format(time.DateTime),
					Utime:   time.UnixMilli(789).Format(time.DateTime),
				},
			},
		},
		{
			name: "get detail of a post of someone else",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					ID:       43,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the article is saved into the DB
				// NOTE: even if this action should fail, we actually got
				// the article, so the cache is updated
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				s.rdb.Del(ctx, "article:content:43")
			},
			param:    "43",
			wantCode: http.StatusBadRequest,
			wantRes: Result[web.ArticleVO]{
				Code: ginx.CodeUserSide,
				Data: web.ArticleVO{},
				Msg:  "article not found",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/articles/detail/%s", tc.param),
				nil,
			)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			rec := httptest.NewRecorder()

			s.server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			var res Result[web.ArticleVO]
			err = json.NewDecoder(rec.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)

			time.Sleep(100 * time.Millisecond) // wait for goroutines
		})
	}
}

func (s *ArticleHandlerSuite) TestList() {
	t := s.T()

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		page web.Page

		wantCode int
		wantRes  Result[[]web.ArticleVO]
	}{
		{
			name: "authors see list of their articles from db",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					ID:       51,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusUnpublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the cache is set
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := s.rdb.Get(ctx, "article:first_page:123").Bytes()
				assert.NoError(t, err)
				var res []domain.Article
				err = json.Unmarshal(val, &res)
				assert.NoError(t, err)
				assert.Equal(t, []domain.Article{{
					ID:      51,
					Title:   "my title",
					Content: "my content",
					Author:  domain.Author{ID: 123},
					Status:  domain.ArticleStatusUnpublished,
					Ctime:   time.UnixMilli(456),
					Utime:   time.UnixMilli(789),
				}}, res)
				assert.NoError(t, s.rdb.Del(ctx, "article:first_page:123").Err())
				_, err = s.rdb.Get(ctx, "article:content:51").Bytes()
				assert.NoError(t, err)
				assert.NoError(t, s.rdb.Del(ctx, "article:content:51").Err())
			},
			page: web.Page{
				Limit:  1,
				Offset: 0,
			},
			wantCode: http.StatusOK,
			wantRes: Result[[]web.ArticleVO]{
				Code: ginx.CodeOK,
				Data: []web.ArticleVO{{
					ID:       51,
					Title:    "my title",
					Abstract: "my content",
					Status:   domain.ArticleStatusUnpublished,
					Ctime:    time.UnixMilli(456).Format(time.DateTime),
					Utime:    time.UnixMilli(789).Format(time.DateTime),
				}},
			},
		},
		{
			name: "authors see list of their articles from cache",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := json.Marshal([]domain.Article{{
					ID:      52,
					Title:   "my title",
					Content: "my content",
					Author:  domain.Author{ID: 123},
					Status:  domain.ArticleStatusUnpublished,
					Ctime:   time.UnixMilli(456),
					Utime:   time.UnixMilli(789),
				}})
				assert.NoError(t, err)
				err = s.rdb.Set(ctx, "article:first_page:123", val, time.Second).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				assert.NoError(t, s.rdb.Del(ctx, "article:first_page:123").Err())
			},
			page: web.Page{
				Limit:  1,
				Offset: 0,
			},
			wantCode: http.StatusOK,
			wantRes: Result[[]web.ArticleVO]{
				Code: ginx.CodeOK,
				Data: []web.ArticleVO{{
					ID:       52,
					Title:    "my title",
					Abstract: "my content",
					Status:   domain.ArticleStatusUnpublished,
					Ctime:    time.UnixMilli(456).Format(time.DateTime),
					Utime:    time.UnixMilli(789).Format(time.DateTime),
				}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(tc.page)
			req, err := http.NewRequest(
				http.MethodPost,
				"/articles/list",
				bytes.NewReader(reqBody),
			)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			rec := httptest.NewRecorder()

			s.server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			var res Result[[]web.ArticleVO]
			err = json.NewDecoder(rec.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)

			time.Sleep(100 * time.Millisecond) // wait for goroutines
		})
	}
}

func (s *ArticleHandlerSuite) TestPubDetail() {
	t := s.T()

	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)

		param string

		wantCode int
		wantRes  Result[web.ArticleVO]
	}{
		{
			name: "readers see detail of an article from db",
			before: func(t *testing.T) {
				err := s.db.Create(dao.PublishedArticle{
					ID:       61,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
				err = s.db.Create(dao.User{
					ID:   123,
					Name: "user",
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// check that the cache is set
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := s.rdb.Get(ctx, "article:published_content:61").Bytes()
				assert.NoError(t, err)
				var res domain.Article
				err = json.Unmarshal(val, &res)
				assert.NoError(t, err)
				assert.Equal(t, domain.Article{
					ID:      61,
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID:   123,
						Name: "user",
					},
					Status: domain.ArticleStatusPublished,
					Ctime:  time.UnixMilli(456),
					Utime:  time.UnixMilli(789),
				}, res)
				assert.NoError(t, s.rdb.Del(ctx, "article:published_content:61").Err())
			},
			param:    "61",
			wantCode: http.StatusOK,
			wantRes: Result[web.ArticleVO]{
				Code: ginx.CodeOK,
				Data: web.ArticleVO{
					ID:         61,
					Title:      "my title",
					Content:    "my content",
					AuthorID:   123,
					AuthorName: "user",
					Status:     domain.ArticleStatusPublished,
					Ctime:      time.UnixMilli(456).Format(time.DateTime),
					Utime:      time.UnixMilli(789).Format(time.DateTime),
				},
			},
		},
		{
			name: "readers see detail of an article from cache",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				val, err := json.Marshal(domain.Article{
					ID:      62,
					Title:   "my title",
					Content: "my content",
					Author: domain.Author{
						ID:   123,
						Name: "user",
					},
					Status: domain.ArticleStatusPublished,
					Ctime:  time.UnixMilli(456),
					Utime:  time.UnixMilli(789),
				})
				assert.NoError(t, err)
				err = s.rdb.Set(ctx, "article:published_content:62", val, time.Second).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				assert.NoError(t, s.rdb.Del(ctx, "article:published_content:62").Err())
			},
			param:    "62",
			wantCode: http.StatusOK,
			wantRes: Result[web.ArticleVO]{
				Code: ginx.CodeOK,
				Data: web.ArticleVO{
					ID:         62,
					Title:      "my title",
					Content:    "my content",
					AuthorID:   123,
					AuthorName: "user",
					Status:     domain.ArticleStatusPublished,
					Ctime:      time.UnixMilli(456).Format(time.DateTime),
					Utime:      time.UnixMilli(789).Format(time.DateTime),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			req, err := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/articles/pub/%s", tc.param),
				nil,
			)
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			rec := httptest.NewRecorder()

			s.server.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantCode, rec.Code)
			var res Result[web.ArticleVO]
			err = json.NewDecoder(rec.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)

			time.Sleep(100 * time.Millisecond) // wait for goroutines
		})
	}
}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}

func init() {
	// limit log output
	gin.SetMode(gin.ReleaseMode)
}
