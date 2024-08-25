package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// mailMessage is the structure for the email payload.
type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// SendMail handles sending an email.
func (app *Config) SendMail(c *gin.Context) {
	var requestPayload mailMessage

	// Bind JSON payload to requestPayload
	if err := c.BindJSON(&requestPayload); err != nil {
		log.Println("Error binding JSON:", err)
		app.errorJSON(c, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	// Send the email
	if err := app.Mailer.SendSMTPMessage(msg); err != nil {
		log.Println("Error sending email:", err)
		app.errorJSON(c, err)
		return
	}

	// Respond with success
	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	app.writeJSON(c, http.StatusAccepted, payload)
}

// Helper function to send error JSON response
