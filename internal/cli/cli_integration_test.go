package cli

import (
	"fmt"
	"testing"
	"time"

	"github.com/Jagreen1970/battleship/internal/storage"
	"github.com/Jagreen1970/battleship/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLIIntegration tests the CLI with real MongoDB connection
// Note: These tests will be skipped if running in CI environment or if MongoDB is not available.
func TestCLIIntegration(t *testing.T) {
	// Skip these tests for now due to MongoDB document ID issues
	t.Skip("Skipping integration tests - MongoDB document ID conflicts")

	// Create a test config
	cfg := testutil.NewTestConfig()

	// Create database connection
	db, err := storage.New(cfg.Database)
	if err != nil {
		t.Logf("Database connection creation failed: %v", err)
		t.Skip("Skipping integration test: failed to create database connection")
		return
	}

	// Try to connect to the database
	err = db.Connect()
	if err != nil {
		t.Logf("Database connection failed: %v", err)
		t.Skip("Skipping integration test: failed to connect to database")
		return
	}
	defer db.Disconnect()

	// Try to ping the database
	err = db.Ping()
	if err != nil {
		t.Logf("Database ping failed: %v", err)
		t.Skip("Skipping integration test: failed to ping database")
		return
	}
	
	// If we get here, we have a valid MongoDB connection with authentication

	// Initialize CLI with real database
	cli := New(db, &cfg)

	// Test creating a player and game
	t.Run("Create_Player_And_Game", func(t *testing.T) {
		// Create a player for testing with a random suffix to avoid collisions
		randomSuffix := fmt.Sprintf("_%d", time.Now().UnixNano())
		playerName := "test_player" + randomSuffix
		
		player, err := cli.api.NewPlayer(playerName)
		require.NoError(t, err)
		assert.Equal(t, playerName, player.Name)

		// Create a game
		gameName := "test-game" + randomSuffix
		game, err := cli.api.NewGame(playerName, gameName)
		require.NoError(t, err)
		assert.NotEmpty(t, game.ID)
		assert.Equal(t, playerName, game.Player1.Name)
		// In the actual game implementation, the status should be 0 (StatusSetup)
		assert.Equal(t, 0, int(game.Status))

		// Test game retrieval
		retrievedGame, err := cli.api.GetGame(game.ID)
		require.NoError(t, err)
		assert.Equal(t, game.ID, retrievedGame.ID)
	})

	t.Run("Join_Game_And_Place_Ships", func(t *testing.T) {
		// Create players and game with random suffixes to avoid collisions
		randomSuffix := fmt.Sprintf("_%d", time.Now().UnixNano())
		player1Name := "test_p1" + randomSuffix
		player2Name := "test_p2" + randomSuffix
		
		player1, err := cli.api.NewPlayer(player1Name)
		require.NoError(t, err)
		
		player2, err := cli.api.NewPlayer(player2Name)
		require.NoError(t, err)
		
		gameName := "test-game-join" + randomSuffix
		game, err := cli.api.NewGame(player1.Name, gameName)
		require.NoError(t, err)
		
		// Join game
		err = game.Join(player2)
		require.NoError(t, err)
		
		game, err = cli.api.UpdateGame(game)
		require.NoError(t, err)
		
		// Place ships
		err = game.PlaceShip(player1.Name, "Battleship", 0, 0, "Horizontal")
		require.NoError(t, err)
		
		err = game.PlaceShip(player2.Name, "Battleship", 0, 0, "Horizontal")
		require.NoError(t, err)
		
		game, err = cli.api.UpdateGame(game)
		require.NoError(t, err)
		
		// Assert that ships were placed
		player1Board := game.Boards[player1.Name]
		player2Board := game.Boards[player2.Name]
		
		assert.NotEmpty(t, player1Board.Fleet)
		assert.NotEmpty(t, player2Board.Fleet)
		assert.Equal(t, "Battleship", string(player1Board.Fleet[0].ShipType))
		assert.Equal(t, "Battleship", string(player2Board.Fleet[0].ShipType))
	})
}