package dbserver

import (
	"gav/internal/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("social.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	DB.AutoMigrate(&user.User{})
}
