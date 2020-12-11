package utils

import (
	"log"
	"os"
	"strconv"
)

// Env get env or fatal if empty!
func Env(envName string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Fatalf("No '%s' env variable", envName)
	}
	return value
}

// EnvAsFloat get env value and convert to float
func EnvAsFloat(envName string, defaultValue float64) float64 {
	var result float64
	envValue := os.Getenv(envName)
	if envValue != "" {
		var err error
		result, err = strconv.ParseFloat(envValue, 64)
		if err != nil {
			log.Fatalf("Could not convert %s %s into float64: %v", envName, envValue, err)
		}
		log.Printf("%s set to %f", envName, result)
	} else {
		result = defaultValue
		log.Printf("%s set to default value of %f", envName, result)
	}
	return result
}

// EnvAsInt get env value and convert to int64
func EnvAsInt(envName string, defaultValue int64) int64 {
	var result int64
	envValue := os.Getenv(envName)
	if envValue != "" {
		var err error
		result, err = strconv.ParseInt(envValue, 10, 64)
		if err != nil {
			log.Fatalf("Could not convert %s %s into int64: %v", envName, envValue, err)
		}
		log.Printf("%s set to %d", envName, result)
	} else {
		result = defaultValue
		log.Printf("%s set to default value of %d", envName, result)
	}
	return result
}
