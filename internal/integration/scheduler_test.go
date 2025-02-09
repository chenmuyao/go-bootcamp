package integration

import (
	"context"
	"testing"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/integration/startup"
	"github.com/chenmuyao/go-bootcamp/internal/job"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type SchedulerTestSuite struct {
	suite.Suite
	scheudler *job.Scheduler
	db        *gorm.DB
}

func (s *SchedulerTestSuite) SetupSuite() {
	s.db = startup.InitDB()
	s.scheudler = startup.InitJobScheduler()
}

func (s *SchedulerTestSuite) TearDownSuite() {
	err := s.db.Exec("TRUNCATE TABLE `jobs`").Error
	assert.NoError(s.T(), err)
}

func (s *SchedulerTestSuite) TestScheduler() {
	testCases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		interval time.Duration
		wantErr  error
		wantJob  *testJob
	}{
		{
			name: "Test Job",
			before: func(t *testing.T) {
				j := dao.Job{
					ID:       1,
					Name:     "test_job",
					Executor: "local",
					// Run every 5 seconds
					CronExpr: "*/5 * * * * ?",
					NextTime: time.Now().UnixMilli(),
					Ctime:    123,
					Utime:    456,
				}
				err := s.db.Create(&j).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var j dao.Job
				err := s.db.Where("id = ?", 1).First(&j).Error
				assert.NoError(t, err)
				assert.True(t, j.NextTime > time.Now().UnixMilli())
				j.NextTime = 0
				assert.True(t, j.Ctime > 0)
				j.Ctime = 0
				assert.True(t, j.Utime > 0)
				j.Utime = 0
				assert.Equal(t, dao.Job{
					ID:       1,
					Name:     "test_job",
					Executor: "local",
					CronExpr: "*/5 * * * * ?",
					Status:   0,
					Version:  1,
				}, j)
			},
			wantErr: context.DeadlineExceeded,
			wantJob: &testJob{cnt: 1},
			// Timeout after 1 second ==> run once
			interval: time.Second,
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			exec := job.NewLocalFuncExecutor()
			j := &testJob{}
			exec.RegisterExecutor("test_job", j.Do)
			s.scheudler.RegisterExecutor(exec)
			ctx, cancel := context.WithTimeout(context.Background(), tc.interval)
			defer cancel()
			err := s.scheudler.Schedule(ctx)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantJob, j)
		})
	}
}

func TestScheduler(t *testing.T) {
	suite.Run(t, &SchedulerTestSuite{})
}

type testJob struct {
	cnt int
}

func (t *testJob) Do(ctx context.Context, j domain.Job) error {
	t.cnt++
	return nil
}
