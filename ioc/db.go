package ioc

import (
	"github.com/chenmuyao/go-bootcamp/config"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(
		mysql.Open(config.Config.DB.DSN),
		&gorm.Config{},
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
