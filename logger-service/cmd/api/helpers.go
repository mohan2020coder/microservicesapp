package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON tries to read the body of a request and converts it into JSON
func (app *Config) readJSON(c *gin.Context, data any) error {
	if err := c.BindJSON(data); err != nil {
		return err
	}
	return nil
}

// writeJSON takes a response status code and arbitrary data and writes a JSON response to the client
func (app *Config) writeJSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

// errorJSON takes an error, and optionally a response status code, and generates and sends
// a JSON error response
func (app *Config) errorJSON(c *gin.Context, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	c.JSON(statusCode, payload)
}
