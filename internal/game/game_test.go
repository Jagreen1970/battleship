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

func TestGameRemoveShip(t *testing.T) {
	type testCase struct {
		name       string
		setupFunc  func(*Game) error
		playerName string
		x, y       int
		wantErr    error
	}

	tests := []testCase{
		{
			name:       "remove ship before any ships placed",
			setupFunc:  func(g *Game) error { return nil },
			playerName: "player1",
			x:          0,
			y:          0,
			wantErr:    ErrorNotFound,
		},
		{
			name: "successfully remove placed ship",
			setupFunc: func(g *Game) error {
				return g.PlaceShip(g.Player1.Name, Submarine, 0, 0, OrientationHorizontal)
			},
			playerName: "player1",
			x:          0,
			y:          0,
			wantErr:    nil,
		},
		{
			name: "remove ship after game started",
			setupFunc: func(g *Game) error {
				if err := g.Join(&Player{Name: "player2"}); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player1.Name, Submarine, 0, 0, OrientationHorizontal); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player2.Name, Submarine, 2, 0, OrientationHorizontal); err != nil {
					return err
				}
				return g.Start(g.Player1.Name)
			},
			playerName: "player1",
			x:          0,
			y:          0,
			wantErr:    ErrorIllegal,
		},
		{
			name:       "remove ship with invalid player",
			setupFunc:  func(g *Game) error { return nil },
			playerName: "invalid_player",
			x:          0,
			y:          0,
			wantErr:    ErrorIllegal,
		},
		{
			name: "remove ship at invalid coordinates",
			setupFunc: func(g *Game) error {
				return g.PlaceShip(g.Player1.Name, Submarine, 0, 0, OrientationHorizontal)
			},
			playerName: "player1",
			x:          5,
			y:          5,
			wantErr:    ErrorIllegal,
		},
		{
			name: "remove ship from opponent's board",
			setupFunc: func(g *Game) error {
				if err := g.Join(&Player{Name: "player2"}); err != nil {
					return err
				}
				return g.PlaceShip("player2", Submarine, 0, 0, OrientationHorizontal)
			},
			playerName: "player1",
			x:          0,
			y:          0,
			wantErr:    ErrorIllegal,
		},
		{
			name: "remove ship during opponent's turn",
			setupFunc: func(g *Game) error {
				if err := g.Join(&Player{Name: "player2"}); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player1.Name, Submarine, 0, 0, OrientationHorizontal); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player2.Name, Submarine, 2, 0, OrientationHorizontal); err != nil {
					return err
				}
				if err := g.Start(g.Player1.Name); err != nil {
					return err
				}
				return g.MakeMove(Move{Player: g.Player1.Name, X: 0, Y: 0})
			},
			playerName: "player1",
			x:          0,
			y:          0,
			wantErr:    ErrorIllegal,
		},
		{
			name: "remove ship at board edge",
			setupFunc: func(g *Game) error {
				return g.PlaceShip(g.Player1.Name, Submarine, 9, 9, OrientationHorizontal)
			},
			playerName: "player1",
			x:          9,
			y:          9,
			wantErr:    nil,
		},
		{
			name: "remove ship with multiple ships on board",
			setupFunc: func(g *Game) error {
				if err := g.PlaceShip(g.Player1.Name, Submarine, 0, 0, OrientationHorizontal); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player1.Name, Cruiser, 2, 0, OrientationHorizontal); err != nil {
					return err
				}
				return g.PlaceShip(g.Player1.Name, Destroyer, 4, 0, OrientationHorizontal)
			},
			playerName: "player1",
			x:          0,
			y:          0,
			wantErr:    nil,
		},
		{
			name: "remove ship from middle of fleet",
			setupFunc: func(g *Game) error {
				if err := g.PlaceShip(g.Player1.Name, Submarine, 0, 0, OrientationHorizontal); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player1.Name, Cruiser, 2, 0, OrientationHorizontal); err != nil {
					return err
				}
				return g.PlaceShip(g.Player1.Name, Destroyer, 4, 0, OrientationHorizontal)
			},
			playerName: "player1",
			x:          2,
			y:          0,
			wantErr:    nil,
		},
		{
			name: "remove ship at coordinates with no ship",
			setupFunc: func(g *Game) error {
				return g.PlaceShip(g.Player1.Name, Submarine, 0, 0, OrientationHorizontal)
			},
			playerName: "player1",
			x:          1,
			y:          1,
			wantErr:    ErrorIllegal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			game := NewGame(&Player{Name: "player1"})
			if err := tc.setupFunc(game); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			// Execute
			err := game.RemoveShip(tc.playerName, tc.x, tc.y)

			// Verify
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				// Verify ship was actually removed if no error expected
				board, exists := game.Boards[tc.playerName]
				if assert.True(t, exists) {
					ships := board.Fleet.Filter(byPosition(tc.x, tc.y))
					assert.Empty(t, ships, "ship should have been removed")
				}
			}
		})
	}
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

func TestGamePlaceShip(t *testing.T) {
	type testCase struct {
		name        string
		setupFunc   func(*Game) error
		playerName  string
		shipType    ShipType
		x, y        int
		orientation ShipOrientation
		wantErr     error
	}

	tests := []testCase{
		{
			name:        "place battleship horizontally",
			setupFunc:   func(g *Game) error { return nil },
			playerName:  "player1",
			shipType:    Battleship,
			x:           0,
			y:           0,
			orientation: OrientationHorizontal,
			wantErr:     nil,
		},
		{
			name:        "place battleship vertically",
			setupFunc:   func(g *Game) error { return nil },
			playerName:  "player1",
			shipType:    Battleship,
			x:           0,
			y:           0,
			orientation: OrientationVertical,
			wantErr:     nil,
		},
		{
			name: "place all ship types",
			setupFunc: func(g *Game) error {
				if err := g.PlaceShip(g.Player1.Name, Battleship, 0, 0, OrientationHorizontal); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player1.Name, Cruiser, 2, 0, OrientationHorizontal); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player1.Name, Destroyer, 4, 0, OrientationHorizontal); err != nil {
					return err
				}
				return g.PlaceShip(g.Player1.Name, Submarine, 6, 0, OrientationHorizontal)
			},
			playerName:  "player1",
			shipType:    Submarine,
			x:           8,
			y:           0,
			orientation: OrientationHorizontal,
			wantErr:     nil,
		},
		{
			name: "place ship after game started",
			setupFunc: func(g *Game) error {
				if err := g.Join(&Player{Name: "player2"}); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player1.Name, Submarine, 0, 0, OrientationHorizontal); err != nil {
					return err
				}
				if err := g.PlaceShip(g.Player2.Name, Submarine, 2, 0, OrientationHorizontal); err != nil {
					return err
				}
				return g.Start(g.Player1.Name)
			},
			playerName:  "player1",
			shipType:    Cruiser,
			x:           4,
			y:           0,
			orientation: OrientationHorizontal,
			wantErr:     ErrorIllegal,
		},
		{
			name:        "place ship out of bounds horizontally",
			setupFunc:   func(g *Game) error { return nil },
			playerName:  "player1",
			shipType:    Battleship,
			x:           8,
			y:           0,
			orientation: OrientationHorizontal,
			wantErr:     ErrorIllegal,
		},
		{
			name:        "place ship out of bounds vertically",
			setupFunc:   func(g *Game) error { return nil },
			playerName:  "player1",
			shipType:    Battleship,
			x:           0,
			y:           8,
			orientation: OrientationVertical,
			wantErr:     ErrorIllegal,
		},
		{
			name: "place overlapping ships",
			setupFunc: func(g *Game) error {
				return g.PlaceShip(g.Player1.Name, Battleship, 0, 0, OrientationHorizontal)
			},
			playerName:  "player1",
			shipType:    Cruiser,
			x:           0,
			y:           0,
			orientation: OrientationVertical,
			wantErr:     ErrorIllegal,
		},
		{
			name: "place ship adjacent horizontally",
			setupFunc: func(g *Game) error {
				return g.PlaceShip(g.Player1.Name, Battleship, 0, 0, OrientationHorizontal)
			},
			playerName:  "player1",
			shipType:    Cruiser,
			x:           0,
			y:           1,
			orientation: OrientationHorizontal,
			wantErr:     ErrorIllegal,
		},
		{
			name: "place ship adjacent vertically",
			setupFunc: func(g *Game) error {
				return g.PlaceShip(g.Player1.Name, Battleship, 0, 0, OrientationVertical)
			},
			playerName:  "player1",
			shipType:    Cruiser,
			x:           1,
			y:           0,
			orientation: OrientationVertical,
			wantErr:     ErrorIllegal,
		},
		{
			name: "place ship diagonally adjacent",
			setupFunc: func(g *Game) error {
				return g.PlaceShip(g.Player1.Name, Battleship, 0, 0, OrientationHorizontal)
			},
			playerName:  "player1",
			shipType:    Cruiser,
			x:           1,
			y:           1,
			orientation: OrientationHorizontal,
			wantErr:     ErrorIllegal,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			game := NewGame(&Player{Name: "player1"})
			if err := tc.setupFunc(game); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			// Execute
			err := game.PlaceShip(tc.playerName, tc.shipType, tc.x, tc.y, tc.orientation)

			// Verify
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
				// Verify ship was actually placed
				board, exists := game.Boards[tc.playerName]
				if assert.True(t, exists) {
					ships := board.Fleet.Filter(func(s *Ship) bool {
						return s.ShipType == tc.shipType &&
							s.IsAtPosition(tc.x, tc.y) &&
							s.Orientation == tc.orientation
					})
					assert.Len(t, ships, 1, "ship should have been placed")
				}
			}
		})
	}
}

func TestGameWinCondition(t *testing.T) {
	player1 := &Player{Name: "player1"}
	player2 := &Player{Name: "player2"}
	game := NewGame(player1)

	// Setup game
	err := game.Join(player2)
	assert.NoError(t, err)

	// Place ships for both players
	err = game.PlaceShip(player1.Name, Submarine, 0, 0, OrientationHorizontal)
	assert.NoError(t, err)
	err = game.PlaceShip(player2.Name, Submarine, 0, 0, OrientationHorizontal)
	assert.NoError(t, err)

	// Start game
	err = game.Start(player1.Name)
	assert.NoError(t, err)

	// Player 1 hits all of Player 2's ships
	moves := []Move{
		{Player: player1.Name, X: 0, Y: 0},
		{Player: player2.Name, X: 5, Y: 5}, // Miss
		{Player: player1.Name, X: 1, Y: 0},
	}

	for _, move := range moves {
		err = game.MakeMove(move)
		assert.NoError(t, err)
	}

	assert.Equal(t, StatusWon, game.Status)
}

func TestGameStateTransitions(t *testing.T) {
	player1 := &Player{Name: "player1"}
	player2 := &Player{Name: "player2"}
	game := NewGame(player1)

	// Initial state
	assert.Equal(t, StatusSetup, game.Status)

	// Join second player
	err := game.Join(player2)
	assert.NoError(t, err)
	assert.Equal(t, StatusSetup, game.Status)

	// Place ships
	err = game.PlaceShip(player1.Name, Submarine, 0, 0, OrientationHorizontal)
	assert.NoError(t, err)
	err = game.PlaceShip(player2.Name, Submarine, 0, 0, OrientationHorizontal)
	assert.NoError(t, err)

	// Start game
	err = game.Start(player1.Name)
	assert.NoError(t, err)
	assert.Equal(t, StatusPlaying, game.Status)

	// Make moves until win
	moves := []Move{
		{Player: player1.Name, X: 0, Y: 0},
		{Player: player2.Name, X: 5, Y: 5},
		{Player: player1.Name, X: 1, Y: 0},
	}

	for i, move := range moves {
		err = game.MakeMove(move)
		assert.NoError(t, err)
		if i < len(moves)-1 {
			assert.Equal(t, StatusPlaying, game.Status)
		}
	}

	assert.Equal(t, StatusWon, game.Status)
}

func TestGameInvalidMoves(t *testing.T) {
	player1 := &Player{Name: "player1"}
	player2 := &Player{Name: "player2"}
	game := NewGame(player1)

	// Setup game
	err := game.Join(player2)
	assert.NoError(t, err)

	err = game.PlaceShip(player1.Name, Submarine, 0, 0, OrientationHorizontal)
	assert.NoError(t, err)
	err = game.PlaceShip(player2.Name, Submarine, 0, 0, OrientationHorizontal)
	assert.NoError(t, err)

	err = game.Start(player1.Name)
	assert.NoError(t, err)

	// Test invalid moves
	invalidMoves := []struct {
		move    Move
		wantErr error
	}{
		{Move{Player: "invalid_player", X: 0, Y: 0}, ErrorIllegal},
		{Move{Player: player2.Name, X: 0, Y: 0}, ErrorIllegal}, // Wrong turn
		{Move{Player: player1.Name, X: -1, Y: 0}, ErrorInvalid},
		{Move{Player: player1.Name, X: 10, Y: 0}, ErrorInvalid},
		{Move{Player: player1.Name, X: 0, Y: -1}, ErrorInvalid},
		{Move{Player: player1.Name, X: 0, Y: 10}, ErrorInvalid},
	}

	for _, tc := range invalidMoves {
		err := game.MakeMove(tc.move)
		assert.Error(t, err)
		assert.ErrorIs(t, err, tc.wantErr)
	}
}

func TestGamePlayerTurns(t *testing.T) {
	player1 := &Player{Name: "player1"}
	player2 := &Player{Name: "player2"}
	game := NewGame(player1)

	// Setup game
	err := game.Join(player2)
	assert.NoError(t, err)

	err = game.PlaceShip(player1.Name, Submarine, 0, 0, OrientationHorizontal)
	assert.NoError(t, err)
	err = game.PlaceShip(player2.Name, Submarine, 0, 0, OrientationHorizontal)
	assert.NoError(t, err)

	err = game.Start(player1.Name)
	assert.NoError(t, err)

	// Verify initial turn
	assert.Equal(t, player1.Name, game.PlayerToMove)

	// Make valid moves and verify turn changes
	moves := []struct {
		move         Move
		nextToMove   string
		expectedHits int
	}{
		{Move{Player: player1.Name, X: 0, Y: 0}, player2.Name, 1},
		{Move{Player: player2.Name, X: 5, Y: 5}, player1.Name, 1},
		{Move{Player: player1.Name, X: 1, Y: 0}, player2.Name, 2},
	}

	for _, tc := range moves {
		err := game.MakeMove(tc.move)
		assert.NoError(t, err)
		assert.Equal(t, tc.nextToMove, game.PlayerToMove)
		assert.Len(t, game.History, tc.expectedHits)
	}
}

// createReadyGame returns a fully initialized game with two players and all ships placed,
// ready to start playing. The game is started with player1's turn.
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
