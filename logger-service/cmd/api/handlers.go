package main

import (
	"logservice/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog handles the logging of data.
func (app *Config) WriteLog(c *gin.Context) {
	// Read JSON into var
	var requestPayload JSONPayload
	if err := c.BindJSON(&requestPayload); err != nil {
		app.errorJSON(c, err)
		return
	}

	// Insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	if err := app.Models.LogEntry.Insert(event); err != nil {
		app.errorJSON(c, err)
		return
	}

	// Respond with success
	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	c.JSON(http.StatusAccepted, resp)
}
