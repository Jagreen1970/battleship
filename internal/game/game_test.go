package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	player1 := &Player{Name: "player1"}
	game := NewGame(player1)

	assert.NotNil(t, game)
	assert.Equal(t, player1, game.Player1)
	assert.Equal(t, "nobody", game.Player2.Name)
	assert.Equal(t, StatusSetup, game.Status)
	assert.Len(t, game.Boards, 1)
	assert.Contains(t, game.Boards, player1.Name)
	assert.Empty(t, game.History)
}

func TestGameJoin(t *testing.T) {
	player1 := &Player{Name: "player1"}
	player2 := &Player{Name: "player2"}
	game := NewGame(player1)

	err := game.Join(player2)
	assert.NoError(t, err)
	assert.Equal(t, player2, game.Player2)
	assert.Equal(t, StatusSetup, game.Status)
	assert.Len(t, game.Boards, 2)
	assert.Contains(t, game.Boards, player1.Name)
	assert.Contains(t, game.Boards, player2.Name)
}

func TestGameJoinError(t *testing.T) {
	player1 := &Player{Name: "player1"}
	player2 := &Player{Name: "player2"}
	game := NewGame(player1)

	err := game.Join(player2)
	assert.NoError(t, err)

	player3 := &Player{Name: "player3"}
	err = game.Join(player3)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrorIllegal)
}

func TestGameOpponent(t *testing.T) {
	player1 := &Player{Name: "player1"}
	game := NewGame(player1)

	// Test opponent before second player joins
	opponent := game.opponent(player1.Name)
	assert.Equal(t, "", opponent)

	// Add second player
	player2 := &Player{Name: "player2"}
	err := game.Join(player2)
	assert.NoError(t, err)

	// Test opponent after second player joins
	opponent = game.opponent(player1.Name)
	assert.Equal(t, player2.Name, opponent)

	opponent = game.opponent(player2.Name)
	assert.Equal(t, player1.Name, opponent)
}

func TestGamePlayerTurns(t *testing.T) {
	game, err := createReadyGame()
	assert.NoError(t, err)

	player1 := game.Player1
	player2 := game.Player2

	// Start game
	err = game.Start(player1.Name)
	assert.NoError(t, err)

	// Verify initial turn
	assert.Equal(t, player1.Name, game.PlayerToMove)

	// Make valid moves and verify turn changes
	moves := []struct {
		move        Move
		nextToMove  string
		histEntries int
	}{
		{Move{Player: player1.Name, X: 0, Y: 0}, player2.Name, 1},
		{Move{Player: player2.Name, X: 5, Y: 5}, player1.Name, 2},
		{Move{Player: player1.Name, X: 1, Y: 0}, player2.Name, 3},
	}

	for _, tc := range moves {
		err := game.MakeMove(tc.move)
		assert.NoError(t, err)
		assert.Equal(t, tc.nextToMove, game.PlayerToMove)
		assert.Len(t, game.History, tc.histEntries)
	}
}

// createReadyGame returns a fully initialized game with two players and all ships placed,
// ready to start playing. The game is not started.
func createReadyGame() (*Game, error) {
	player1 := &Player{Name: "player1"}
	player2 := &Player{Name: "player2"}
	game := NewGame(player1)

	// Join second player
	err := game.Join(player2)
	if err != nil {
		return nil, err
	}

	// Place ships for both players
	ships := []struct {
		playerName  string
		shipType    ShipType
		x, y        int
		orientation ShipOrientation
	}{
		// Player 1 ships (horizontal layout)
		{player1.Name, Battleship, 0, 0, OrientationHorizontal}, // 1 Battleship
		{player1.Name, Cruiser, 0, 2, OrientationHorizontal},    // 2 Cruisers
		{player1.Name, Cruiser, 5, 2, OrientationHorizontal},
		{player1.Name, Destroyer, 0, 4, OrientationHorizontal}, // 3 Destroyers
		{player1.Name, Destroyer, 4, 4, OrientationHorizontal},
		{player1.Name, Destroyer, 0, 6, OrientationHorizontal},
		{player1.Name, Submarine, 4, 6, OrientationHorizontal}, // 4 Submarines
		{player1.Name, Submarine, 8, 6, OrientationHorizontal},
		{player1.Name, Submarine, 0, 8, OrientationHorizontal},
		{player1.Name, Submarine, 3, 8, OrientationHorizontal},

		// Player 2 ships (vertical layout)
		{player2.Name, Battleship, 0, 0, OrientationHorizontal}, // 1 Battleship
		{player2.Name, Cruiser, 0, 2, OrientationHorizontal},    // 2 Cruisers
		{player2.Name, Cruiser, 5, 2, OrientationHorizontal},
		{player2.Name, Destroyer, 0, 4, OrientationHorizontal}, // 3 Destroyers
		{player2.Name, Destroyer, 4, 4, OrientationHorizontal},
		{player2.Name, Destroyer, 0, 6, OrientationHorizontal},
		{player2.Name, Submarine, 4, 6, OrientationHorizontal}, // 4 Submarines
		{player2.Name, Submarine, 8, 6, OrientationHorizontal},
		{player2.Name, Submarine, 0, 8, OrientationHorizontal},
		{player2.Name, Submarine, 3, 8, OrientationHorizontal},
	}

	for _, ship := range ships {
		err = game.PlaceShip(ship.playerName, ship.shipType, ship.x, ship.y, ship.orientation)
		if err != nil {
			return nil, err
		}
	}

	return game, nil
}

// TestGameplaySetup tests the setup of a game with two players and all ships placed.
func TestGameplaySetup(t *testing.T) {
	game, err := createReadyGame()
	assert.NoError(t, err)

	// Start the game with player1's turn
	err = game.Start(game.Player1.Name)
	assert.NoError(t, err)

	assert.Equal(t, StatusPlaying, game.Status)
	assert.Equal(t, "player1", game.PlayerToMove)

	// Verify correct number of ships for each player
	for _, playerName := range []string{"player1", "player2"} {
		board := game.Boards[playerName]
		assert.Len(t, board.Fleet, FleetSizeAllowed)

		// Verify correct number of each ship type
		for shipType, expectedCount := range shipsAllowed {
			ships := board.Fleet.Filter(byShipType(shipType))
			assert.Len(t, ships, expectedCount, "Wrong number of %v for %s", shipType, playerName)
		}
	}
}
