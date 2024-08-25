package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON tries to read the body of a request and converts it into JSON
func (app *Config) ReadJSON(c *gin.Context, data any) error {
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
func (app *Config) writeJSON(c *gin.Context, status int, data any, headers ...map[string][]string) error {
	_, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, values := range headers[0] {
			for _, value := range values {
				c.Header(key, value)
			}
		}
	}

	c.Header("Content-Type", "application/json")
	c.JSON(status, data)

	return nil
}

func (app *Config) errorJSON(c *gin.Context, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}

	return app.writeJSON(c, statusCode, payload)
}
