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
	s.db.Exec("TRUNCATE user_like_bizs")
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
					Biz:     "read",
					BizID:   1,
					ReadCnt: 1,
				}, intr)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HGet(ctx, "interactive:read:1", "read_cnt").Err()
				assert.ErrorIs(t, err, redis.Nil)
				err = s.rdb.Del(ctx, "interactive:read:1").Err()
				assert.NoError(t, err)
			},
			bizID: 1,
		},
		{
			name: "read count incr when cache is absent",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Interactive{
					ID:      12,
					Biz:     "read",
					BizID:   2,
					ReadCnt: 2,
					Ctime:   123,
					Utime:   123,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				err := s.db.Where("id = ?", 12).First(&intr).Error
				t.Log(intr)
				assert.NoError(t, err)
				assert.True(t, intr.Utime > 123)
				intr.Utime = 0
				assert.Equal(t, dao.Interactive{
					ID:      12,
					Biz:     "read",
					BizID:   2,
					ReadCnt: 3,
					Ctime:   123,
				}, intr)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HGet(ctx, "interactive:read:22", "read_cnt").Err()
				assert.ErrorIs(t, err, redis.Nil)
				err = s.rdb.Del(ctx, "interactive:read:22").Err()
				assert.NoError(t, err)
			},
			bizID: 2,
		},
		{
			name: "increment read count from db and redis",
			before: func(t *testing.T) {
				// set readCnt in DB
				err := s.db.Create(dao.Interactive{
					ID:      13,
					Biz:     "read",
					BizID:   3,
					ReadCnt: 3,
					Ctime:   123,
					Utime:   123,
				}).Error
				assert.NoError(t, err)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HSet(ctx, "interactive:read:3", "read_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				err := s.db.Where("id = ?", 13).First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Utime > 123)
				intr.Utime = 0
				assert.Equal(t, dao.Interactive{
					ID:      13,
					Biz:     "read",
					BizID:   3,
					ReadCnt: 4,
					Ctime:   123,
				}, intr)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := s.rdb.HGet(ctx, "interactive:read:3", "read_cnt").Int()
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
			err := svc.IncrReadCnt(context.Background(), "read", tc.bizID)
			assert.Equal(t, nil, err)
		})
	}
}

func (s *InteractivceTestSuite) TestLike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		bizID  int64
	}{
		{
			name: "like count incr when both db and cache empty",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				err := s.db.Where("biz_id = ? AND biz = ?", 1, "like").First(&intr).Error
				assert.NoError(t, err)

				assert.True(t, intr.ID > 0)
				intr.ID = 0
				assert.True(t, intr.Ctime > 0)
				intr.Ctime = 0
				assert.True(t, intr.Utime > 0)
				intr.Utime = 0
				assert.Equal(t, dao.Interactive{
					Biz:     "like",
					BizID:   1,
					LikeCnt: 1,
				}, intr)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("uid = ? AND biz_id = ? AND biz = ?", 123, 1, "like").
					First(&likeBiz).
					Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.ID > 0)
				likeBiz.ID = 0
				assert.True(t, likeBiz.Ctime > 0)
				likeBiz.Ctime = 0
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserLikeBiz{
					UID:    123,
					Biz:    "like",
					BizID:  1,
					Status: 1,
				}, likeBiz)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HGet(ctx, "interactive:like:1", "like_cnt").Err()
				assert.ErrorIs(t, err, redis.Nil)
				err = s.rdb.Del(ctx, "interactive:like:1").Err()
				assert.NoError(t, err)
			},
			bizID: 1,
		},
		{
			name: "like count incr when cache is absent",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Interactive{
					ID:      22,
					Biz:     "like",
					BizID:   2,
					LikeCnt: 2,
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
					Biz:     "like",
					BizID:   2,
					LikeCnt: 3,
					Ctime:   123,
				}, intr)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("uid = ? AND biz_id = ? AND biz = ?", 123, 2, "like").
					First(&likeBiz).
					Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.ID > 0)
				likeBiz.ID = 0
				assert.True(t, likeBiz.Ctime > 0)
				likeBiz.Ctime = 0
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserLikeBiz{
					UID:    123,
					Biz:    "like",
					BizID:  2,
					Status: 1,
				}, likeBiz)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HGet(ctx, "interactive:like:22", "like_cnt").Err()
				assert.ErrorIs(t, err, redis.Nil)
				err = s.rdb.Del(ctx, "interactive:like:22").Err()
				assert.NoError(t, err)
			},
			bizID: 2,
		},
		{
			name: "increment like count from db and redis",
			before: func(t *testing.T) {
				// set likeCnt in DB
				err := s.db.Create(dao.Interactive{
					ID:      23,
					Biz:     "like",
					BizID:   3,
					LikeCnt: 3,
					Ctime:   123,
					Utime:   123,
				}).Error
				assert.NoError(t, err)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HSet(ctx, "interactive:like:3", "like_cnt", 3).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				err := s.db.Where("id = ?", 23).First(&intr).Error
				assert.NoError(t, err)
				assert.True(t, intr.Utime > 123)
				intr.Utime = 0
				assert.Equal(t, dao.Interactive{
					ID:      23,
					Biz:     "like",
					BizID:   3,
					LikeCnt: 4,
					Ctime:   123,
				}, intr)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("uid = ? AND biz_id = ? AND biz = ?", 123, 3, "like").
					First(&likeBiz).
					Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.ID > 0)
				likeBiz.ID = 0
				assert.True(t, likeBiz.Ctime > 0)
				likeBiz.Ctime = 0
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserLikeBiz{
					UID:    123,
					Biz:    "like",
					BizID:  3,
					Status: 1,
				}, likeBiz)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := s.rdb.HGet(ctx, "interactive:like:3", "like_cnt").Int()
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
			err := svc.Like(context.Background(), "like", tc.bizID, 123)
			assert.Equal(t, nil, err)
		})
	}
}

func (s *InteractivceTestSuite) TestCancelLike() {
	t := s.T()
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		bizID  int64
	}{
		{
			// We don't expected any error, nothing should happen.
			name: "decr like count when both db and cache empty",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				err := s.db.Where("biz = ?", "cancel_like").First(&intr).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
				var likeBiz dao.UserLikeBiz
				err = s.db.First(&likeBiz).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
			bizID: 1,
		},
		{
			name: "decr like count incr when cache is absent",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Interactive{
					ID:      32,
					Biz:     "cancel_like",
					BizID:   2,
					LikeCnt: 2,
					Ctime:   123,
					Utime:   123,
				}).Error
				assert.NoError(t, err)
				err = s.db.Create(dao.UserLikeBiz{
					ID:     32,
					UID:    123,
					Biz:    "cancel_like",
					BizID:  2,
					Status: 0,
					Ctime:  123,
					Utime:  123,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var intr dao.Interactive
				err := s.db.Where("id = ?", 32).First(&intr).Error
				t.Log(intr)
				assert.NoError(t, err)
				assert.True(t, intr.Utime > 123)
				intr.Utime = 0
				assert.Equal(t, dao.Interactive{
					ID:      32,
					Biz:     "cancel_like",
					BizID:   2,
					LikeCnt: 1,
					Ctime:   123,
				}, intr)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("uid = ? AND biz_id = ? AND biz = ?", 123, 2, "cancel_like").
					First(&likeBiz).
					Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.ID > 0)
				likeBiz.ID = 0
				assert.True(t, likeBiz.Ctime > 0)
				likeBiz.Ctime = 0
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserLikeBiz{
					UID:    123,
					Biz:    "cancel_like",
					BizID:  2,
					Status: 0,
				}, likeBiz)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HGet(ctx, "interactive:cancel_like:22", "like_cnt").Err()
				assert.ErrorIs(t, err, redis.Nil)
				err = s.rdb.Del(ctx, "interactive:cancel_like:22").Err()
				assert.NoError(t, err)
			},
			bizID: 2,
		},
		{
			name: "decr like count from db and redis",
			before: func(t *testing.T) {
				// set likeCnt in DB
				err := s.db.Create(dao.Interactive{
					ID:      33,
					Biz:     "cancel_like",
					BizID:   3,
					LikeCnt: 3,
					Ctime:   123,
					Utime:   123,
				}).Error
				assert.NoError(t, err)
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				err = s.rdb.HSet(ctx, "interactive:cancel_like:3", "like_cnt", 3).Err()
				assert.NoError(t, err)
				err = s.db.Create(dao.UserLikeBiz{
					ID:     33,
					UID:    123,
					Biz:    "cancel_like",
					BizID:  3,
					Status: 1,
					Ctime:  123,
					Utime:  123,
				}).Error
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
					Biz:     "cancel_like",
					BizID:   3,
					LikeCnt: 2,
					Ctime:   123,
				}, intr)

				var likeBiz dao.UserLikeBiz
				err = s.db.Where("uid = ? AND biz_id = ? AND biz = ?", 123, 3, "cancel_like").
					First(&likeBiz).
					Error
				assert.NoError(t, err)
				assert.True(t, likeBiz.ID > 0)
				likeBiz.ID = 0
				assert.True(t, likeBiz.Ctime > 0)
				likeBiz.Ctime = 0
				assert.True(t, likeBiz.Utime > 0)
				likeBiz.Utime = 0
				assert.Equal(t, dao.UserLikeBiz{
					UID:    123,
					Biz:    "cancel_like",
					BizID:  3,
					Status: 0,
				}, likeBiz)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				res, err := s.rdb.HGet(ctx, "interactive:cancel_like:3", "like_cnt").Int()
				assert.NoError(t, err)
				assert.Equal(t, 2, res)
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
			err := svc.CancelLike(context.Background(), "cancel_like", tc.bizID, 123)
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
