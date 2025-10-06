package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"mysocial/user"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("social.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB.AutoMigrate(&user.User{})
}
