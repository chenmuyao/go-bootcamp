package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	s.db.Exec("truncate table `articles`")
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
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)
				assert.True(t, article.ID > 0)
				assert.Equal(t, "my title", article.Title)
				assert.Equal(t, "my content", article.Content)
			},
			article: web.ArticleEditReq{
				Title:   "my title",
				Content: "my content",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: ginx.CodeOK,
				Data: 1,
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
