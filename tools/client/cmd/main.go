package main

import (
	"log"
	"spacecraftsim/tools/client"
)

func main() {
	// Create a new client
	cli, err := client.New("localhost:8080")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer cli.Close()

	// Start the CLI interface
	cli.RunCLI()
}
