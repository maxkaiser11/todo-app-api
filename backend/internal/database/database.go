package database

import (
	"os"

	"github.com/maxkaiser11/todo-app-api/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDBPath() string {
	if v := os.Getenv("DB_PATH"); v != "" {
		return v
	}
	return "test.db"
}

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(getDBPath()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.User{}, &models.Task{})
	return db, nil
}
