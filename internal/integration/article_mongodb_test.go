package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/integration/startup"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/net/context"
)

type MongoDBArticleHandlerSuite struct {
	suite.Suite
	client   *mongo.Client
	coll     *mongo.Collection
	liveColl *mongo.Collection
	server   *gin.Engine
}

func (s *MongoDBArticleHandlerSuite) SetupSuite() {
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("user", ijwt.UserClaims{
			UID: 123,
		})
	})
	s.client = startup.InitMongoDB()
	// s.rdb = startup.InitRedis()
	s.coll = s.client.Database(dao.DatabaseName).Collection(dao.ArticleCollName)
	s.liveColl = s.client.Database(dao.DatabaseName).Collection(dao.PublishedArticleCollName)
	node, err := snowflake.NewNode(1)
	assert.NoError(s.T(), err)
	hdl := startup.InitArticleHandler(dao.NewMongoDBArticleDAO(s.client, node))
	hdl.RegisterRoutes(s.server)
}

func (s *MongoDBArticleHandlerSuite) TearDownTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := s.coll.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
	s.liveColl.DeleteMany(ctx, bson.D{})
	assert.NoError(s.T(), err)
}

func (s *MongoDBArticleHandlerSuite) TestEdit() {
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
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				filter := bson.M{
					"author_id": 123,
				}
				var article dao.Article
				err := s.coll.FindOne(ctx, filter).Decode(&article)
				assert.NoError(t, err)

				assert.True(t, article.ID > 0)
				article.ID = 0
				assert.True(t, article.Ctime > 789)
				article.Ctime = 0
				assert.True(t, article.Utime > 789)
				article.Utime = 0
				assert.Equal(t, dao.Article{
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
				Data: 1,
				Code: ginx.CodeOK,
			},
		},
		{
			name: "Edit a post",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				article := dao.Article{
					ID:       2,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 123,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}

				_, err := s.coll.InsertOne(ctx, &article)
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				// check that the article is saved into the DB
				var article dao.Article
				err := s.coll.FindOne(ctx, bson.M{
					"id": 2,
				}).Decode(&article)
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
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				article := dao.Article{
					ID:       3,
					Title:    "my title",
					Content:  "my content",
					AuthorID: 234,
					Status:   domain.ArticleStatusPublished,
					Ctime:    456,
					Utime:    789,
				}

				_, err := s.coll.InsertOne(ctx, &article)
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				// check that the article is saved into the DB
				var article dao.Article
				err := s.coll.FindOne(ctx, bson.M{
					"id": 3,
				}).Decode(&article)
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
			if tc.wantRes.Data > 0 {
				// can only assert the existence of an ID
				assert.True(t, res.Data > 0)
				tc.wantRes.Data = res.Data
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

// NOTE: Need to set replica set for mongo
// func (s *MongoDBArticleHandlerSuite) TestPublish() {
// 	t := s.T()

// 	testCases := []struct {
// 		name   string
// 		before func(t *testing.T)
// 		after  func(t *testing.T)

// 		// json article from frontend
// 		article web.ArticlePublishReq

// 		wantCode int
// 		wantRes  Result[int64]
// 	}{
// 		{
// 			name:   "Publish a new post",
// 			before: func(t *testing.T) {},
// 			after: func(t *testing.T) {
// 				// check that the article is saved into the DB
// 				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 				defer cancel()

// 				filter := bson.M{
// 					"author_id": 123,
// 				}
// 				var article dao.Article
// 				err := s.liveColl.FindOne(ctx, filter).Decode(&article)
// 				assert.NoError(t, err)

// 				assert.True(t, article.ID > 0)
// 				article.ID = 0
// 				assert.True(t, article.Ctime > 789)
// 				article.Ctime = 0
// 				assert.True(t, article.Utime > 789)
// 				article.Utime = 0
// 				assert.Equal(t, dao.Article{
// 					Title:    "my super title",
// 					Content:  "my content",
// 					AuthorID: 123,
// 					Status:   domain.ArticleStatusPublished,
// 				}, article)
// 			},
// 			article: web.ArticlePublishReq{
// 				Title:   "my super title",
// 				Content: "my content",
// 			},
// 			wantCode: http.StatusOK,
// 			wantRes: Result[int64]{
// 				Code: ginx.CodeOK,
// 				Data: 1,
// 			},
// 		},
// 		// {
// 		// 	name: "Edit and publish a post",
// 		// 	before: func(t *testing.T) {
// 		// 		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 		// 		defer cancel()
// 		//
// 		// 		article := dao.Article{
// 		// 			ID:       2,
// 		// 			Title:    "my title",
// 		// 			Content:  "my content",
// 		// 			AuthorID: 123,
// 		// 			Status:   domain.ArticleStatusPublished,
// 		// 			Ctime:    456,
// 		// 			Utime:    789,
// 		// 		}
// 		//
// 		// 		_, err := s.liveColl.InsertOne(ctx, &article)
// 		// 		assert.NoError(t, err)
// 		// 	},
// 		// 	after: func(t *testing.T) {
// 		// 		// check that the article is saved into the DB
// 		// 		var article dao.Article
// 		// 		err := s.mdb.Where("id = ?", 2).First(&article).Error
// 		// 		assert.NoError(t, err)
// 		// 		assert.True(t, article.Utime > 789)
// 		// 		article.Utime = 0
// 		// 		assert.Equal(t, dao.Article{
// 		// 			ID:       2,
// 		// 			Title:    "new title",
// 		// 			Content:  "new content",
// 		// 			AuthorID: 123,
// 		// 			Status:   domain.ArticleStatusPublished,
// 		// 			Ctime:    456,
// 		// 		}, article)
// 		// 	},
// 		// 	article: web.ArticlePublishReq{
// 		// 		ID:      2,
// 		// 		Title:   "new title",
// 		// 		Content: "new content",
// 		// 	},
// 		// 	wantCode: http.StatusOK,
// 		// 	wantRes: Result[int64]{
// 		// 		Code: ginx.CodeOK,
// 		// 		Data: 2,
// 		// 	},
// 		// },
// 		// {
// 		// 	name: "Publish a post of someone else",
// 		// 	before: func(t *testing.T) {
// 		// 		err := s.mdb.Create(dao.Article{
// 		// 			ID:       3,
// 		// 			Title:    "my title",
// 		// 			Content:  "my content",
// 		// 			AuthorID: 234,
// 		// 			Status:   domain.ArticleStatusPublished,
// 		// 			Ctime:    456,
// 		// 			Utime:    789,
// 		// 		}).Error
// 		// 		assert.NoError(t, err)
// 		// 	},
// 		// 	after: func(t *testing.T) {
// 		// 		// check that the article is saved into the DB
// 		// 		var article dao.Article
// 		// 		err := s.mdb.Where("id = ?", 3).First(&article).Error
// 		// 		assert.NoError(t, err)
// 		// 		assert.Equal(t, dao.Article{
// 		// 			ID:       3,
// 		// 			Title:    "my title",
// 		// 			Content:  "my content",
// 		// 			AuthorID: 234,
// 		// 			Status:   domain.ArticleStatusPublished,
// 		// 			Ctime:    456,
// 		// 			Utime:    789,
// 		// 		}, article)
// 		// 	},
// 		// 	article: web.ArticlePublishReq{
// 		// 		ID:      3,
// 		// 		Title:   "new title",
// 		// 		Content: "new content",
// 		// 	},
// 		// 	wantCode: http.StatusBadRequest,
// 		// 	wantRes: Result[int64]{
// 		// 		Code: ginx.CodeUserSide,
// 		// 		Msg:  "article not found",
// 		// 	},
// 		// },
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.before(t)
// 			defer tc.after(t)

// 			reqBody, err := json.Marshal(tc.article)
// 			assert.NoError(t, err)
// 			req, err := http.NewRequest(
// 				http.MethodPost,
// 				"/articles/publish",
// 				bytes.NewReader(reqBody),
// 			)
// 			req.Header.Set("Content-Type", "application/json")
// 			assert.NoError(t, err)
// 			rec := httptest.NewRecorder()

// 			s.server.ServeHTTP(rec, req)

// 			assert.Equal(t, tc.wantCode, rec.Code)
// 			var res Result[int64]
// 			err = json.NewDecoder(rec.Body).Decode(&res)
// 			assert.NoError(t, err)
// 			if tc.wantRes.Data > 0 {
// 				// can only assert the existence of an ID
// 				assert.True(t, res.Data > 0)
// 				tc.wantRes.Data = res.Data
// 			}
// 			assert.Equal(t, tc.wantRes, res)
// 		})
// 	}
// }

func TestMongoDBArticleHandler(t *testing.T) {
	suite.Run(t, &MongoDBArticleHandlerSuite{})
}

func init() {
	// limit log output
	gin.SetMode(gin.ReleaseMode)
}
