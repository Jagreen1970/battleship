package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/cli"
	"github.com/Jagreen1970/battleship/internal/server"
	"github.com/Jagreen1970/battleship/internal/storage"
)

func main() {
	// Parse command line flags
	cliMode := flag.Bool("cli", false, "Start in CLI mode")
	dbUser := flag.String("dbuser", "", "MongoDB username")
	dbPass := flag.String("dbpass", "", "MongoDB password") 
	flag.Parse()

	// Load configuration
	cfg, err := app.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set MongoDB credentials from command line if provided
	if *dbUser != "" {
		cfg.Database.User = *dbUser
	}
	if *dbPass != "" {
		cfg.Database.Password = *dbPass
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize database
	db, err := storage.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Connect to database
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Verify connection with ping
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	defer func() {
		if err := db.Disconnect(); err != nil {
			log.Fatalf("Failed to disconnect from database: %v", err)
		}
	}()

	if *cliMode {
		runCLIMode(cfg, db)
		return
	}

	// Initialize server
	s := server.New(cfg.Server)

	// Start server
	go func() {
		if err := s.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	log.Println("Shutting down...")
	if err := s.Shutdown(); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
}

func runCLIMode(cfg *app.Config, db storage.Storage) {
	c := cli.New(db, cfg)
	c.Run()
}
