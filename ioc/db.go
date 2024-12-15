package ioc

import (
	"github.com/chenmuyao/go-bootcamp/config"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.Logger) *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Cfg.DB.DSN), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			SlowThreshold: 0,
			LogLevel:      glogger.Info,
		}),
	},
	)
	if err != nil {
		panic("failed to connect database")
	}

	// TODO: Replace by sql migration
	err = dao.InitTable(db)
	if err != nil {
		panic("failed to init tables")
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(s string, i ...interface{}) {
	g(s, logger.Field{
		Key:   "args",
		Value: i,
	})
}
