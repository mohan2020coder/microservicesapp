package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// routes sets up Gin routes and middleware
func (app *Config) routes() *gin.Engine {
	router := gin.Default()

	// Set up CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Update with your allowed origins in production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Heartbeat endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Define routes
	router.POST("/", app.Broker)
	router.POST("/log-grpc", app.LogViaGRPC)
	router.POST("/handle", app.HandleSubmission)

	return router
}
