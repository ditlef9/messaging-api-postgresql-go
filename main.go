// File: main.go

package main

import (
	"ekeberg.com/messaging-api-postgresql-go/db"
	"ekeberg.com/messaging-api-postgresql-go/handlers"
	"ekeberg.com/messaging-api-postgresql-go/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	// SQLite connection
	db.InitDB()

	// Start Gin Router
	r := gin.Default()

	// Use the CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve static files for the favicon
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Users (no authentication required)
		v1.POST("users/signup", handlers.SignUp) // POST http://localhost:8080/api/v1/users/signup
		v1.POST("users/login", handlers.Login)   // POST http://localhost:8080/api/v1/users/login

		// Messages (authentication required as service)
		authenticatedHumanOrService := v1.Group("/")
		authenticatedHumanOrService.Use(middlewares.Authenticate)
		{
			authenticatedHumanOrService.GET("messages", handlers.GetMessages)        // GET http://localhost:8080/api/v1/messages
			authenticatedHumanOrService.GET("messages/:id", handlers.GetMessageById) // GET http://localhost:8080/api/v1/messages/1
		}

	}

	// By default it serves on :8080 unless a PORT environment variable was defined.
	r.Run()
}
