package dao

import (
	"context"
	"database/sql"
	"math"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserInsert(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(t *testing.T) *sql.DB
		ctx      context.Context
		user     User
		wantErr  error
		wantUser User
	}{
		{
			name: "insert ok",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(123, 123)
				mock.ExpectExec("INSERT INTO .*").
					WithArgs().
					WillReturnError(nil).
					WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			user: User{
				Name: "test",
			},
			wantErr: nil,
			wantUser: User{
				ID:   123,
				Name: "test",
			},
		},
		{
			name: "email conflict",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx: context.Background(),
			user: User{
				Name: "conflict@test.com",
			},
			wantErr:  ErrDuplicatedUser,
			wantUser: User{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: sqlDB,
				// mock has no version
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// mock doesn't respond to ping
				DisableAutomaticPing: true,
				// use plein commands, don't add bigin-commit automatically
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)

			dao := NewUserDAO(db)

			u, err := dao.Insert(tc.ctx, tc.user)
			if err == nil {
				assert.True(
					t,
					math.Abs(float64(u.Utime)-float64(time.Now().UnixMilli())) < float64(1000),
				)
				assert.True(
					t,
					math.Abs(float64(u.Ctime)-float64(time.Now().UnixMilli())) < float64(1000),
				)
				u.Utime = 0
				u.Ctime = 0
			}
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}
