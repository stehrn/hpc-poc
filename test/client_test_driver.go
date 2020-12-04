package main

import (
	"fmt"
	"log"

	"github.com/stehrn/hpc-poc/client"
)

//
func main() {
	log.Print("Starting test driver")

	business := client.BusinessFromEnv()
	client := client.NewEnvClientOrFatal()

	var n int
	for n < 5 {
		data := []byte(fmt.Sprintf("payload %d", n))
		location, id, err := client.Handle(business, data)
		if err != nil {
			log.Fatalf("client.handle() err: %v", err)
		}

		log.Printf("Payload uploaded to cloud storage location: %s, notification sent with message ID: %s", location, id)
	}
}
