package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// routes sets up the Gin router and middleware.
func (app *Config) routes() *gin.Engine {
	// Create a new Gin router
	r := gin.Default()

	// Configure CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Configure heartbeat endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")
	})

	// Define routes
	r.POST("/send", app.SendMail)

	return r
}
