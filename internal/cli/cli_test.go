package cli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/Jagreen1970/battleship/internal/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLI represents the main test for CLI functionality
func TestCLI(t *testing.T) {
	// Create mock database and test config
	mockDB := newMockStorage(t)
	cfg := &app.Config{
		Database: app.DatabaseConfig{
			Name: "test",
		},
	}

	// Initialize CLI with mock database
	cli := New(mockDB, cfg)

	// Test basic CLI functionality
	t.Run("CLI_Initialization", func(t *testing.T) {
		assert.NotNil(t, cli)
		assert.NotNil(t, cli.api)
		assert.NotNil(t, cli.reader)
		assert.Empty(t, cli.currentGameID)
	})

	t.Run("Create_Game", func(t *testing.T) {
		// Setup
		player := "testPlayer"
		gameID := "game123"
		gameName := "Test Game"

		// Add player to mock
		mockDB.players[player] = &game.Player{Name: player}

		// Make createGame set a test game ID
		mockDB.mockCreateGame = func(g *game.Game) (*game.Game, error) {
			g.ID = gameID
			g.Name = gameName
			mockDB.games[gameID] = g
			return g, nil
		}

		// Execute
		cli.createGame(player, gameName)

		// Assert
		assert.Equal(t, gameID, cli.currentGameID)
		assert.Equal(t, gameName, mockDB.games[gameID].Name)
	})

	t.Run("Join_Game", func(t *testing.T) {
		// Setup
		gameID := "game123"
		player1 := "testPlayer"
		player2 := "otherPlayer"

		// Create test game with player1
		testGame := game.NewGame(&game.Player{Name: player1}, "")
		testGame.ID = gameID
		mockDB.games[gameID] = testGame

		// Add player2 to mock
		mockDB.players[player2] = &game.Player{Name: player2}

		// Execute
		cli.joinGame(gameID, player2)

		// Assert
		assert.Contains(t, mockDB.games[gameID].Boards, player2)
		assert.Equal(t, player2, mockDB.games[gameID].Player2.Name)
	})

	t.Run("Set_Game", func(t *testing.T) {
		// Setup
		gameID := "game456"
		testGame := game.NewGame(&game.Player{Name: "player1"}, "")
		testGame.ID = gameID
		mockDB.games[gameID] = testGame

		// Execute
		cli.setGame(gameID)

		// Assert
		assert.Equal(t, gameID, cli.currentGameID)
	})

	t.Run("Place_Ship", func(t *testing.T) {
		// Setup
		gameID := "game123"
		player := "testPlayer"
		shipType := "Battleship"
		x, y := 0, 0
		orientation := "Horizontal"

		// Create a clean game for this test
		testGame := game.NewGame(&game.Player{Name: player}, "")
		testGame.ID = gameID
		mockDB.games[gameID] = testGame

		// Execute
		cli.placeShip(gameID, player, shipType, x, y, orientation)

		// Assert
		board := mockDB.games[gameID].Boards[player]
		assert.NotNil(t, board)
		assert.NotEmpty(t, board.Fleet)
		assert.Equal(t, game.ShipType(shipType), board.Fleet[0].ShipType)
	})

	t.Run("Fire_Shot", func(t *testing.T) {
		// Setup - Create a ready game with two players and ships placed
		gameID := "game789"
		player1 := "player1"
		player2 := "player2"

		// Create and prepare a game for testing
		testGame, err := createTestGameWithShips(player1, player2)
		require.NoError(t, err)
		testGame.ID = gameID

		// Start the game
		err = testGame.Start(player1)
		require.NoError(t, err)

		mockDB.games[gameID] = testGame
		cli.currentGameID = gameID

		// Execute - player1 fires at player2's board
		cli.fire(gameID, player1, 0, 0)

		// Assert
		assert.Len(t, mockDB.games[gameID].History, 1)
		assert.Equal(t, player1, mockDB.games[gameID].History[0].Player)
		assert.Equal(t, 0, mockDB.games[gameID].History[0].X)
		assert.Equal(t, 0, mockDB.games[gameID].History[0].Y)
	})
}

// TestCLI_HandleCommand tests the command handling function
func TestCLI_HandleCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		setup   func(*mockStorage)
		assert  func(*testing.T, *CLI, *mockStorage)
	}{
		{
			name:    "Create Game",
			command: "create-game player1",
			setup: func(m *mockStorage) {
				m.players["player1"] = &game.Player{Name: "player1"}
				m.mockCreateGame = func(g *game.Game) (*game.Game, error) {
					g.ID = "game123"
					return g, nil
				}
			},
			assert: func(t *testing.T, c *CLI, m *mockStorage) {
				assert.Equal(t, "game123", c.currentGameID)
			},
		},
		{
			name:    "Create Game With Name",
			command: "create-game player1 myAwesomeGame",
			setup: func(m *mockStorage) {
				m.players["player1"] = &game.Player{Name: "player1"}
				m.mockCreateGame = func(g *game.Game) (*game.Game, error) {
					g.ID = "game456"
					return g, nil
				}
			},
			assert: func(t *testing.T, c *CLI, m *mockStorage) {
				assert.Equal(t, "game456", c.currentGameID)
				assert.Equal(t, "myAwesomeGame", m.games["game456"].Name)
			},
		},
		{
			name:    "Join Game",
			command: "join-game game123 player2",
			setup: func(m *mockStorage) {
				// Create game and player
				g := game.NewGame(&game.Player{Name: "player1"}, "")
				g.ID = "game123"
				m.games["game123"] = g
				m.players["player2"] = &game.Player{Name: "player2"}
			},
			assert: func(t *testing.T, c *CLI, m *mockStorage) {
				assert.Equal(t, "game123", c.currentGameID)
				assert.Contains(t, m.games["game123"].Boards, "player2")
			},
		},
		{
			name:    "Set Game",
			command: "set-game game123",
			setup: func(m *mockStorage) {
				g := game.NewGame(&game.Player{Name: "player1"}, "")
				g.ID = "game123"
				m.games["game123"] = g
			},
			assert: func(t *testing.T, c *CLI, m *mockStorage) {
				assert.Equal(t, "game123", c.currentGameID)
			},
		},
		{
			name:    "Show Games",
			command: "show-games 0 2",
			setup: func(m *mockStorage) {
				// Create several games for pagination testing
				for i := 0; i < 5; i++ {
					g := game.NewGame(&game.Player{Name: fmt.Sprintf("player%d", i+1)}, "")
					g.ID = fmt.Sprintf("game%d", i+1)
					m.games[g.ID] = g
				}
			},
			assert: func(t *testing.T, c *CLI, m *mockStorage) {
				// Nothing to verify on the CLI state, just ensure no panics
				// Output testing is done in a separate test
			},
		},
		{
			name:    "Delete Game",
			command: "delete-game game123",
			setup: func(m *mockStorage) {
				g := game.NewGame(&game.Player{Name: "player1"}, "")
				g.ID = "game123"
				m.games["game123"] = g
			},
			assert: func(t *testing.T, c *CLI, m *mockStorage) {
				// Verify the game was deleted
				_, exists := m.games["game123"]
				assert.False(t, exists)
			},
		},
		{
			name:    "Place Ship",
			command: "place-ship player1 Battleship 0 0 Horizontal",
			setup: func(m *mockStorage) {
				g := game.NewGame(&game.Player{Name: "player1"}, "")
				g.ID = "game123"
				m.games["game123"] = g
			},
			assert: func(t *testing.T, c *CLI, m *mockStorage) {
				board := m.games["game123"].Boards["player1"]
				assert.NotEmpty(t, board.Fleet)
				assert.Equal(t, game.ShipType("Battleship"), board.Fleet[0].ShipType)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock and CLI
			mockDB := newMockStorage(t)
			cfg := &app.Config{
				Database: app.DatabaseConfig{
					Name: "test",
				},
			}
			cli := New(mockDB, cfg)

			// Additional setup for the test case
			if tc.setup != nil {
				tc.setup(mockDB)
			}

			// Save current game ID (for restoration in specific tests)
			cli.currentGameID = "game123"

			// Execute the command
			cli.handleCommand(tc.command)

			// Assert
			if tc.assert != nil {
				tc.assert(t, cli, mockDB)
			}
		})
	}
}

// TestCLI_InputOutput tests the CLI input/output functionality
func TestCLI_InputOutput(t *testing.T) {
	// Setup mock and CLI with captured output
	mockDB := newMockStorage(t)
	cfg := &app.Config{
		Database: app.DatabaseConfig{
			Name: "test",
		},
	}

	// Create a CLI with string input
	input := "create-game player1\nexit\n"
	reader := strings.NewReader(input)

	// Create a CLI with captured output
	var outputBuffer bytes.Buffer

	cli := New(mockDB, cfg)
	cli.SetIO(reader, &outputBuffer)

	// Setup mock for create-game command
	mockDB.players["player1"] = &game.Player{Name: "player1"}
	mockDB.mockCreateGame = func(g *game.Game) (*game.Game, error) {
		g.ID = "game123"
		return g, nil
	}

	// Instead of running the full Run method (which enters a loop and waits for user input),
	// we'll test the command handling directly
	cli.currentGameID = ""
	cli.handleCommand("create-game player1")

	// Verify the output contains the expected message
	output := outputBuffer.String()
	assert.Contains(t, output, "Created new game with ID: game123")

	// Verify the game ID was set
	assert.Equal(t, "game123", cli.currentGameID)
}

// TestShowGames tests the show-games output functionality
func TestShowGames(t *testing.T) {
	// Setup mock and CLI with captured output
	mockDB := newMockStorage(t)
	cfg := &app.Config{
		Database: app.DatabaseConfig{
			Name: "test",
		},
	}

	// Create a CLI with captured output
	var outputBuffer bytes.Buffer

	cli := New(mockDB, cfg)
	cli.SetIO(nil, &outputBuffer)

	// Create several test games with different statuses
	setupGames := []struct {
		id       string
		player1  string
		player2  string
		status   game.Status
		hasMoves bool
	}{
		{"game1", "player1", "", game.StatusSetup, false},
		{"game2", "player1", "player2", game.StatusPlaying, true},
		{"game3", "player3", "player4", game.StatusWon, true},
		{"game4", "player5", "player6", game.StatusLost, true},
	}

	for _, setup := range setupGames {
		player1 := &game.Player{Name: setup.player1}
		g := game.NewGame(player1, "")
		g.ID = setup.id

		if setup.player2 != "" {
			player2 := &game.Player{Name: setup.player2}
			err := g.Join(player2)
			require.NoError(t, err)
		}

		// Set game status
		g.Status = setup.status

		// Add some moves if needed
		if setup.hasMoves {
			g.History = append(g.History, game.Move{
				Player: setup.player1,
				X:      0,
				Y:      0,
				Hit:    true,
			})
		}

		mockDB.games[setup.id] = g
	}

	// Test pagination with page 0, count 2
	outputBuffer.Reset()
	cli.showGames(0, 2)

	// Verify the output contains the expected pagination controls
	output := outputBuffer.String()
	assert.Contains(t, output, "For next page: show-games 1 2")
	// Don't assert specific game IDs since map iteration order is non-deterministic

	// Test pagination with page 1, count 2
	outputBuffer.Reset()
	cli.showGames(1, 2)

	// Verify the output contains the expected pagination controls
	output = outputBuffer.String()
	assert.Contains(t, output, "For previous page: show-games 0 2")
	// Don't assert specific game IDs since map iteration order is non-deterministic

	// Test empty page
	outputBuffer.Reset()
	cli.showGames(5, 2)

	// Verify the output indicates no games
	output = outputBuffer.String()
	assert.Contains(t, output, "No games found on page 5")
}

// TestDeleteGame tests the delete-game functionality
func TestDeleteGame(t *testing.T) {
	// Setup
	mockDB := newMockStorage(t)
	cfg := &app.Config{
		Database: app.DatabaseConfig{
			Name: "test",
		},
	}

	// Create a CLI with captured output
	var outputBuffer bytes.Buffer
	cli := New(mockDB, cfg)
	cli.SetIO(nil, &outputBuffer)

	// Create a test game
	gameID := "game123"
	player := &game.Player{Name: "player1"}
	g := game.NewGame(player, "")
	g.ID = gameID
	mockDB.games[gameID] = g

	// Set the current game ID
	cli.currentGameID = gameID

	// Test deleting the game
	cli.deleteGame(gameID)

	// Verify the game was deleted
	_, err := mockDB.FindGameByID(gameID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, game.ErrorNotFound)

	// Verify the current game ID was cleared
	assert.Empty(t, cli.currentGameID)

	// Verify the output
	output := outputBuffer.String()
	assert.Contains(t, output, "Game game123 successfully deleted")
	assert.Contains(t, output, "Current game unset")

	// Test deleting a non-existent game
	outputBuffer.Reset()
	cli.deleteGame("nonexistent")

	// Verify the output
	output = outputBuffer.String()
	assert.Contains(t, output, "Error: Game nonexistent not found")
}

// TestDeleteAllGames tests the delete-all-games functionality
func TestDeleteAllGames(t *testing.T) {
	// Setup
	mockDB := newMockStorage(t)
	cfg := &app.Config{
		Database: app.DatabaseConfig{
			Name: "test",
		},
	}

	// Create several test games
	for i := 0; i < 3; i++ {
		g := game.NewGame(&game.Player{Name: fmt.Sprintf("player%d", i+1)}, "")
		g.ID = fmt.Sprintf("game%d", i+1)
		mockDB.games[g.ID] = g
	}

	// Set up a CLI with simulated input/output
	confirmInput := "confirm\n"
	inputReader := strings.NewReader(confirmInput)
	var outputBuffer bytes.Buffer

	cli := New(mockDB, cfg)
	cli.SetIO(inputReader, &outputBuffer)
	cli.currentGameID = "game1"

	// Test deleting all games
	cli.deleteAllGames()

	// Verify all games were deleted
	assert.Empty(t, mockDB.games)

	// Verify the current game ID was cleared
	assert.Empty(t, cli.currentGameID)

	// Verify the output
	output := outputBuffer.String()
	assert.Contains(t, output, "WARNING: This will delete ALL games")
	assert.Contains(t, output, "Successfully deleted 3 games")
	assert.Contains(t, output, "Current game unset")

	// Test cancellation
	// Reset the mock and CLI
	mockDB = newMockStorage(t)
	for i := 0; i < 3; i++ {
		g := game.NewGame(&game.Player{Name: fmt.Sprintf("player%d", i+1)}, "")
		g.ID = fmt.Sprintf("game%d", i+1)
		mockDB.games[g.ID] = g
	}

	cancelInput := "cancel\n"
	inputReader = strings.NewReader(cancelInput)
	outputBuffer.Reset()

	cli = New(mockDB, cfg)
	cli.SetIO(inputReader, &outputBuffer)
	cli.currentGameID = "game1"

	// Test cancelling the deletion
	cli.deleteAllGames()

	// Verify no games were deleted
	assert.Len(t, mockDB.games, 3)

	// Verify the current game ID was not cleared
	assert.Equal(t, "game1", cli.currentGameID)

	// Verify the output
	output = outputBuffer.String()
	assert.Contains(t, output, "WARNING: This will delete ALL games")
	assert.Contains(t, output, "Operation cancelled")
}

// createTestGameWithShips is a helper function to create a test game with ships placed
func createTestGameWithShips(player1Name, player2Name string) (*game.Game, error) {
	player1 := &game.Player{Name: player1Name}
	player2 := &game.Player{Name: player2Name}

	g := game.NewGame(player1, "")

	// Join second player
	err := g.Join(player2)
	if err != nil {
		return nil, err
	}

	// Place enough ships to deplete pins and make the game ready
	ships := []struct {
		playerName  string
		shipType    string
		x, y        int
		orientation string
	}{
		// Player 1 ships
		{player1Name, "Battleship", 0, 0, "Horizontal"},
		{player1Name, "Cruiser", 0, 2, "Horizontal"},
		{player1Name, "Cruiser", 5, 2, "Horizontal"},
		{player1Name, "Destroyer", 0, 4, "Horizontal"},
		{player1Name, "Destroyer", 4, 4, "Horizontal"},
		{player1Name, "Destroyer", 0, 6, "Horizontal"},
		{player1Name, "Submarine", 4, 6, "Horizontal"},
		{player1Name, "Submarine", 8, 6, "Horizontal"},
		{player1Name, "Submarine", 0, 8, "Horizontal"},
		{player1Name, "Submarine", 3, 8, "Horizontal"},

		// Player 2 ships
		{player2Name, "Battleship", 0, 0, "Horizontal"},
		{player2Name, "Cruiser", 0, 2, "Horizontal"},
		{player2Name, "Cruiser", 5, 2, "Horizontal"},
		{player2Name, "Destroyer", 0, 4, "Horizontal"},
		{player2Name, "Destroyer", 4, 4, "Horizontal"},
		{player2Name, "Destroyer", 0, 6, "Horizontal"},
		{player2Name, "Submarine", 4, 6, "Horizontal"},
		{player2Name, "Submarine", 8, 6, "Horizontal"},
		{player2Name, "Submarine", 0, 8, "Horizontal"},
		{player2Name, "Submarine", 3, 8, "Horizontal"},
	}

	for _, ship := range ships {
		err = g.PlaceShip(ship.playerName, game.ShipType(ship.shipType), ship.x, ship.y, game.ShipOrientation(ship.orientation))
		if err != nil {
			// In tests sometimes placing all ships might not be possible due to space constraints
			// That's fine as long as we have placed enough ships to start the game
			continue
		}
	}

	// Set pins available to 0 to force game to be ready to start
	for _, board := range g.Boards {
		board.PinsAvailable = 0
	}

	return g, nil
}
