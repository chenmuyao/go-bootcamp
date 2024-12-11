package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	// NOTE: Not the best practice. Too risky. Strong dependency
	return db.AutoMigrate(&User{}, &SMSInfo{})
}
