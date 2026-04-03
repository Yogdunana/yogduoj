package main

import (
	"log"

	"github.com/Yogdunana/yogduoj/judge/internal/config"
	"github.com/Yogdunana/yogduoj/judge/internal/server"
)

func main() {
	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting YogduOJ Judge Service...")

	// Create and initialize the server (container pool, workers, gRPC).
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Handle graceful shutdown in a separate goroutine.
	go srv.WaitForShutdown()

	// Start serving (blocking).
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Judge service stopped.")
}
