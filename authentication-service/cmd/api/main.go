package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"authentication/data"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *gorm.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// Connect to DB
	conn, err := connectToDB()
	if err != nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToDB() (*gorm.DB, error) {
	dsn := os.Getenv("DSN")
	for {
		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
			if counts > 10 {
				return nil, err
			}
			log.Println("Backing off for two seconds....")
			time.Sleep(2 * time.Second)
			continue
		}
		log.Println("Connected to Postgres!")
		return conn, nil
	}
}
