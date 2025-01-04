package integration

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/integration/startup"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache/rediscache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type InteractivceTestSuite struct {
	suite.Suite
	db  *gorm.DB
	rdb redis.Cmdable
}

func (s *InteractivceTestSuite) SetupSuite() {
	s.db = startup.InitDB()
	s.rdb = startup.InitRedis()
}

func (s *InteractivceTestSuite) TearDownTest() {
	log.Println("TRUNCATE")
	s.db.Exec("TRUNCATE interactives")
}

func (s *InteractivceTestSuite) TestIncrReadCnt() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		bizID  int64
	}{
		{
			name: "read count incr when both db and cache empty",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				err := s.db.Where("id = ?", 1).First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Ctime > 0)
				intr.Ctime = 0
				assert.True(t, intr.Utime > 0)
				intr.Utime = 0
				assert.Equal(t, dao.Interactive{
					ID:      1,
					Biz:     "test",
					BizID:   1,
					ReadCnt: 1,
				}, intr)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HGet(ctx, "interactive:test:1", "read_cnt").Err()
				assert.ErrorIs(t, err, redis.Nil)
				err = s.rdb.Del(ctx, "interactive:test:1").Err()
				assert.NoError(t, err)
			},
			bizID: 1,
		},
		{
			name: "read count incr when cache is absent",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Interactive{
					ID:      22,
					Biz:     "test",
					BizID:   2,
					ReadCnt: 2,
					Ctime:   123,
					Utime:   123,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				err := s.db.Where("id = ?", 22).First(&intr).Error
				t.Log(intr)
				assert.NoError(t, err)
				assert.True(t, intr.Utime > 123)
				intr.Utime = 0
				assert.Equal(t, dao.Interactive{
					ID:      22,
					Biz:     "test",
					BizID:   2,
					ReadCnt: 3,
					Ctime:   123,
				}, intr)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HGet(ctx, "interactive:test:22", "read_cnt").Err()
				assert.ErrorIs(t, err, redis.Nil)
				err = s.rdb.Del(ctx, "interactive:test:22").Err()
				assert.NoError(t, err)
			},
			bizID: 2,
		},
		{
			name: "increment read count from db and redis",
			before: func(t *testing.T) {
				// set readCnt in DB
				err := s.db.Create(dao.Interactive{
					ID:      33,
					Biz:     "test",
					BizID:   3,
					ReadCnt: 3,
					Ctime:   123,
					Utime:   123,
				}).Error
				assert.NoError(t, err)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HSet(ctx, "interactive:test:3", "read_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				err := s.db.Where("id = ?", 33).First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Utime > 123)
				intr.Utime = 0
				assert.Equal(t, dao.Interactive{
					ID:      33,
					Biz:     "test",
					BizID:   3,
					ReadCnt: 4,
					Ctime:   123,
				}, intr)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := s.rdb.HGet(ctx, "interactive:test:3", "read_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 4, res)
			},
			bizID: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			dao := dao.NewGORMInteractiveDAO(s.db)
			cache := rediscache.NewInteractiveRedisCache(s.rdb)
			repo := repository.NewCachedInteractiveRepository(dao, cache)
			svc := service.NewInteractiveService(repo)
			err := svc.IncrReadCnt(context.Background(), "test", tc.bizID)
			assert.Equal(t, nil, err)
		})
	}
}

func TestInteractiveSvc(t *testing.T) {
	suite.Run(t, &InteractivceTestSuite{})
}

func init() {
	// limit log output
	gin.SetMode(gin.ReleaseMode)
}
