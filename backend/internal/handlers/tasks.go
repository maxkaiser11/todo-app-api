package handlers

import (
	"net/http"

	"github.com/maxkaiser11/todo-app-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		var tasks []models.Task
		db.Where("user_id = ?", userID).Find(&tasks)
		c.JSON(http.StatusOK, tasks)
	}
}

func GetTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		id := c.Param("id")
		var t models.Task
		if err := db.Where("id = ? AND user_id = ?", id, userID).First(&t).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusOK, t)
	}
}

func CreateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		var t models.Task
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		t.UserID = userID
		db.Create(&t)
		c.JSON(http.StatusCreated, t)
	}
}

func UpdateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		id := c.Param("id")
		var t models.Task
		if err := db.Where("id = ? AND user_id = ?", id, userID).First(&t).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		var body models.Task
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		t.Status = body.Status
		db.Save(&t)
		c.JSON(http.StatusOK, t)
	}
}

func DeleteTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		id := c.Param("id")
		if err := db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	}
}
