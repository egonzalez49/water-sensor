package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/egonzalez49/water-sensor/subscriber"
	"github.com/joho/godotenv"
)

func main() {
	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, os.Interrupt, syscall.SIGTERM)

	loadEnvVars()

	subscriber.Subscribe()

	<-keepAlive
}

func loadEnvVars() {
	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file.")
	}
}
