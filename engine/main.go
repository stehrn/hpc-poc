package main

import (
	"log"
	"os"
)

func main() {
	log.Print("Starting Engine")
	payload := env("PAYLOAD")
	log.Printf("Engine got payload %s", payload)
}

func env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("No '%s' env variable", key)
	}
	return value
}
