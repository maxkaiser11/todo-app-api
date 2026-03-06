package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Task        string
	Description string
	Status      bool
}

func main() {
	db, err := gorm.Open(sqlite.Open("/data/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connnect database")
	}

	db.AutoMigrate(&Task{})

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "https://todo-app-api-frontend.vercel.app"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT"},
		AllowHeaders: []string{"Content-Type"},
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "welcome to my API",
		})
	})

	router.GET("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		var t Task
		if err := db.First(&t, "id == ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusOK, t)
	})

	router.POST("/task/create", func(c *gin.Context) {
		var t Task
		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "task created", "task": t})
		}
		addTask(t, db)
		c.JSON(http.StatusCreated, gin.H{"message": "task created", "task": t})
	})

	router.GET("/tasks", func(c *gin.Context) {
		var tasks []Task
		db.Find(&tasks)
		c.JSON(http.StatusOK, tasks)
	})

	router.DELETE("/task/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		if err := db.Delete(&Task{}, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	})

	fmt.Println("Server running on Port 8080 🚀")
	router.Run("0.0.0.0:8080")
}
