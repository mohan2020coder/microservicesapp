package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Config) Authenticate(c *gin.Context) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&requestPayload); err != nil {
		app.errorJSON(c, err, http.StatusBadRequest)
		return
	}

	// Validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(c, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(c, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// Log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	payload := gin.H{
		"error":   false,
		"message": fmt.Sprintf("Logged in user %s", user.Email),
		"data":    user,
	}

	c.JSON(http.StatusAccepted, payload)
}
func (app *Config) logRequest(name, data string) error {
	entry := struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}{
		Name: name,
		Data: data,
	}

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		return err
	}

	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("error logging request: status code %d", response.StatusCode)
	}

	return nil
}
