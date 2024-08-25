package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"broker/event"
	"broker/logs"
)

// RequestPayload describes the JSON that this service accepts as an HTTP Post request
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

// MailPayload is the embedded type (in RequestPayload) that describes an email message to be sent
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// AuthPayload is the embedded type (in RequestPayload) that describes an authentication request
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LogPayload is the embedded type (in RequestPayload) that describes a request to log something
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Broker is a test handler, just to make sure we can hit the broker from a web client
func (app *Config) Broker(c *gin.Context) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	app.writeJSON(c, http.StatusOK, payload)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(c *gin.Context) {
	var requestPayload RequestPayload

	err := app.readJSON(c, &requestPayload)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(c, requestPayload.Auth)
	case "log":
		app.logItemViaRPC(c, requestPayload.Log)
	case "mail":
		app.sendMail(c, requestPayload.Mail)
	default:
		app.errorJSON(c, errors.New("unknown action"))
	}
}

// logItem logs an item by making an HTTP Post request with a JSON payload, to the logger microservice
func (app *Config) LogItem(c *gin.Context, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(c, errors.New("error logging item"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(c, http.StatusAccepted, payload)
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(c *gin.Context, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(c, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(c, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(c, errors.New("authentication error"), http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(c, http.StatusAccepted, payload)
}

// sendMail sends email by calling the mail microservice
func (app *Config) sendMail(c *gin.Context, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceURL := "http://mailer-service/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(c, errors.New("error calling mail service"))
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Message sent to " + msg.To

	app.writeJSON(c, http.StatusAccepted, payload)
}

// logEventViaRabbit logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (app *Config) LogEventViaRabbit(c *gin.Context, l LogPayload) {
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via RabbitMQ"

	app.writeJSON(c, http.StatusAccepted, payload)
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}
	l := LogPayload{
		Name: name,
		Data: msg,
	}

	payload := ConvertLogPayloadToRPCPayload(l)

	j, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

// logItemViaRPC logs an item by making an RPC call to the logger microservice
func (app *Config) logItemViaRPC(c *gin.Context, l LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer client.Close()

	// Convert LogPayload to RPCPayload using the conversion function
	rpcPayload := ConvertLogPayloadToRPCPayload(l)

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: result,
	}

	app.writeJSON(c, http.StatusAccepted, payload)
}

// ConvertLogPayloadToRPCPayload converts a LogPayload to RPCPayload
func ConvertLogPayloadToRPCPayload(logPayload LogPayload) RPCPayload {
	return RPCPayload{
		Name: logPayload.Name,
		Data: logPayload.Data,
	}
}

// LogViaGRPC handles the gRPC call to write logs
func (app *Config) LogViaGRPC(c *gin.Context) {
	var requestPayload RequestPayload

	// Read JSON payload from the request
	err := app.readJSON(c, &requestPayload)
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	// Convert LogPayload to RPCPayload using the conversion function
	rpcPayload := ConvertLogPayloadToRPCPayload(requestPayload.Log)

	// Create a new gRPC client connection using NewClient
	conn, err := grpc.NewClient("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		app.errorJSON(c, err)
		return
	}
	defer conn.Close()

	// Create a new gRPC client
	logClient := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Make the gRPC call to WriteLog
	_, err = logClient.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: rpcPayload.Name,
			Data: rpcPayload.Data,
		},
	})
	if err != nil {
		app.errorJSON(c, err)
		return
	}

	// Respond to the client
	payload := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(c, http.StatusAccepted, payload)
}
