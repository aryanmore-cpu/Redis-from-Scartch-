package main

import (
	"log"
)

func main() {
	server, err := NewServer("127.0.0.1:6379")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
