package startup

import (
	"github.com/chenmuyao/go-bootcamp/interactive/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(
		mysql.Open(
			"root:root@tcp(127.0.0.1:13316)/wetravel?charset=utf8mb4&parseTime=True&loc=Local",
		),
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
