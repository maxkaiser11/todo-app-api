package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string `json:"-"`
}

type Task struct {
	gorm.Model
	Task        string `json:"Task"`
	Description string `json:"Description"`
	Status      bool   `json:"Status"`
	UserID      uint   `json:"UserID"`
}
