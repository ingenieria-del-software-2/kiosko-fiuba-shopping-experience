package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/api"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/config"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := database.NewPostgresConnectionWithConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize API server
	srv := api.NewServer(db, cfg)

	// Start server
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	log.Printf("Shopping Experience Microservice started on %s:%d", cfg.Host, cfg.Port)

	// Wait for termination signal
	select {
	case <-sigCh:
		log.Println("Shutdown signal received")
	case <-ctx.Done():
		log.Println("Server error occurred")
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server shutdown gracefully")
}
