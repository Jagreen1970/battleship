package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/game"
	"github.com/Jagreen1970/battleship/internal/storage"
)

func main() {
	// Create a test configuration with MongoDB credentials
	cfg := app.Config{
		Database: app.DatabaseConfig{
			Driver:   "mongo",
			URL:      "mongodb://localhost:27017",
			Name:     "battleship",
			Timeout:  5 * time.Second,
			User:     "root",
			Password: "battleship",
		},
	}

	// Create and connect to the database
	db, err := storage.New(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}

	err = db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Disconnect()

	// Create the game API
	api := game.NewApi(db)

	// Test player creation with unique names
	testPlayer1 := fmt.Sprintf("test_player_%d", time.Now().UnixNano())
	player1, err := api.NewPlayer(testPlayer1)
	if err != nil {
		log.Fatalf("Failed to create player1: %v", err)
	}
	fmt.Printf("Created player1: %s (ID: %s)\n", player1.Name, player1.ID)

	// Try to create another player with a different name
	testPlayer2 := fmt.Sprintf("test_player_%d", time.Now().UnixNano())
	player2, err := api.NewPlayer(testPlayer2)
	if err != nil {
		log.Fatalf("Failed to create player2: %v", err)
	}
	fmt.Printf("Created player2: %s (ID: %s)\n", player2.Name, player2.ID)

	// Create a game
	game, err := api.NewGame(player1.Name)
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}
	fmt.Printf("Created game: %s\n", game.ID)

	// Join the game with player2
	err = game.Join(player2)
	if err != nil {
		log.Fatalf("Failed to join game: %v", err)
	}
	fmt.Printf("Player %s joined the game\n", player2.Name)

	// Update the game
	updatedGame, err := api.UpdateGame(game)
	if err != nil {
		log.Fatalf("Failed to update game: %v", err)
	}
	fmt.Printf("Game updated. Status: %d\n", updatedGame.Status)

	fmt.Println("Test completed successfully!")
}