package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connnect database")
	}
	return db
}
