package main

import (
	"log"
	"os"
)

// A simple engine, well, a bit of code that reads and ENV variable and prints it out.
// Alternative and more realistic approach would load file(s) from a mouted filesystem

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
