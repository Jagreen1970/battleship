package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/game"
	"github.com/Jagreen1970/battleship/internal/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	newGame, err := api.NewGame(player1.Name)
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}
	fmt.Printf("Created game with ID: %s\n", newGame.ID)
	fmt.Printf("Game ID length: %d\n", len(newGame.ID))
	
	// Verify if the ID is a valid ObjectID
	_, err = primitive.ObjectIDFromHex(newGame.ID)
	if err != nil {
		fmt.Printf("ERROR: Game ID is not a valid ObjectID: %v\n", err)
	} else {
		fmt.Println("Game ID is a valid ObjectID")
	}

	// Retrieve the game to verify
	retrievedGame, err := api.GetGame(newGame.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve game: %v", err)
	}
	fmt.Printf("Retrieved game with ID: %s\n", retrievedGame.ID)

	// Join the game with player2
	err = retrievedGame.Join(player2)
	if err != nil {
		log.Fatalf("Failed to join game: %v", err)
	}
	fmt.Printf("Player %s joined the game\n", player2.Name)

	// Debug - print game ID before update
	fmt.Printf("Game ID before update: %s\n", retrievedGame.ID)
	fmt.Printf("Game ID length before update: %d\n", len(retrievedGame.ID))

	// Update the game
	updatedGame, err := api.UpdateGame(retrievedGame)
	if err != nil {
		log.Fatalf("Failed to update game: %v\n", err)
	}
	fmt.Printf("Game updated. Status: %d\n", updatedGame.Status)

	fmt.Println("Test completed successfully!")
}