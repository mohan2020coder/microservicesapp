package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Println("Failed to create consumer:", err)
		os.Exit(1)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println("Failed to start listening to the queue:", err)
		os.Exit(1)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	rabbitMQURL := "amqp://guest:guest@rabbitmq" // This can be set from an environment variable if needed.

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready, retrying...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("Failed to connect to RabbitMQ after multiple attempts:", err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off for", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
