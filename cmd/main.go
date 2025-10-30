package main

import (
	"fmt"
	"log"

	"github.com/Indroneel007/Load-Balancer/internal/server"
)

func main() {
	fmt.Println("Starting Load Balancer Server...")
	err := server.Run()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
