package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/api"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/database"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Initialize database connection
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize API server
	srv := api.NewServer(db)

	// Start server
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	log.Println("Shopping Experience Microservice started")

	// Wait for termination signal
	select {
	case <-sigCh:
		log.Println("Shutdown signal received")
	case <-ctx.Done():
		log.Println("Server error occurred")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), server.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server shutdown gracefully")
}
