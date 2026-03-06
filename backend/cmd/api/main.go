package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/maxkaiser11/todo-app-api/internal/database"
	"github.com/maxkaiser11/todo-app-api/internal/handlers"
	"github.com/maxkaiser11/todo-app-api/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	godotenv.Load()
	db, err := database.Connect()
	if err != nil {
		panic("failed to connect database")
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "https://todo-app-api-frontend.vercel.app"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	// Public
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "welcome to my API"})
	})
	router.POST("/auth/register", handlers.Register(db))
	router.POST("/auth/login", handlers.Login(db))

	// Protected
	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.GET("/tasks", handlers.GetTasks(db))
	auth.GET("/tasks/:id", handlers.GetTask(db))
	auth.POST("/task/create", handlers.CreateTask(db))
	auth.PUT("/tasks/:id", handlers.UpdateTask(db))
	auth.DELETE("/task/delete/:id", handlers.DeleteTask(db))

	fmt.Println("Server running on Port 8080 🚀")
	router.Run("0.0.0.0:8080")
}
