package utils

import (
	"log"
	"os"
)

// Env get env or fatal if empty!
func Env(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("No '%s' env variable", key)
	}
	return value
}
