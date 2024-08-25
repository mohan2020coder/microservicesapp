package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// readJSON tries to read the body of a request and converts it into JSON
func (app *Config) ReadJSON(c *gin.Context, data interface{}) error {
	maxBytes := 1048576 // one megabyte

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(maxBytes))

	dec := json.NewDecoder(c.Request.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON takes a response status code and arbitrary data and writes a JSON response to the client
// writeJSON takes a response status code and arbitrary data and writes a JSON response to the client
func (app *Config) writeJSON(c *gin.Context, status int, data interface{}, headers ...map[string]string) {
	out, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			c.Header(key, value)
		}
	}

	c.Data(status, "application/json", out)
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
