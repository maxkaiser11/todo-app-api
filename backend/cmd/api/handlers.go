package main

import "gorm.io/gorm"

func addTask(task Task, db *gorm.DB) {
	db.Create(&Task{Task: task.Task, Description: task.Description, Status: task.Status})
}
