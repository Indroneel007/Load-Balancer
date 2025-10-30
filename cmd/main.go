package main

import (
	"log"

	"github.com/Indroneel007/Load-Balancer/internal/server"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
