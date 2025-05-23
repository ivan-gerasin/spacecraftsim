package main

import (
	"log"
	"spacecraftsim/internal/server"
)

func main() {
	// Create and start the server
	srv := server.New(":8080")
	log.Printf("Starting TCP server on :8080...")

	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
