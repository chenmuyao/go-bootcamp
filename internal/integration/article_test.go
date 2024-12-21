package integration

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/integration/startup"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/gin-gonic/gin"
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
	server *gin.Engine
}

func (s *ArticleHandlerSuite) SetupSuite() {
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UserClaims{
			UID: 123,
		})
	})
	hdl := startup.InitArticleHandler()
	hdl.RegisterRoutes(s.server)
	s.db = startup.InitDB()
	// s.rdb = startup.InitRedis()
}

func (s *ArticleHandlerSuite) TearDownTest() {
	log.Println("TRUNCATE")
	s.db.Exec("TRUNCATE articles")
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

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}

func init() {
	// limit log output
	gin.SetMode(gin.ReleaseMode)
}
