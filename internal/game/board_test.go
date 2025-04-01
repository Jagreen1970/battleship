package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
	// Test data
	playerName := "player1"
	opponentName := "player2"

	// Create new board
	board := NewBoard(playerName, opponentName)

	// Assert board is not nil
	assert.NotNil(t, board)

	// Check initial pins available
	assert.Equal(t, 30, board.PinsAvailable)

	// Check maps initialization
	assert.Len(t, board.Maps, 2)
	assert.NotNil(t, board.Maps[0])
	assert.NotNil(t, board.Maps[1])

	// Check map titles
	assert.Equal(t, playerName, board.Maps[0].Title)
	assert.Equal(t, opponentName, board.Maps[1].Title)

	// Check that fleet is initially nil
	assert.Nil(t, board.Fleet)

	// Check that all fields are initialized to FieldStateEmpty
	for x := range 10 {
		for y := range 10 {
			assert.Equal(t, FieldStateEmpty, board.Maps[0].FieldState(x, y),
				"Player map field at (%d,%d) should be empty", x, y)
			assert.Equal(t, FieldStateEmpty, board.Maps[1].FieldState(x, y),
				"Opponent map field at (%d,%d) should be empty", x, y)
		}
	}
}

func TestBoardPlaceShip(t *testing.T) {
	tests := []struct {
		name          string
		shipType      ShipType
		x, y          int
		orientation   ShipOrientation
		expectedError bool
	}{
		{
			name:          "Place battleship horizontally",
			shipType:      Battleship,
			x:             0,
			y:             0,
			orientation:   OrientationHorizontal,
			expectedError: false,
		},
		{
			name:          "Place cruiser vertically",
			shipType:      Cruiser,
			x:             9,
			y:             0,
			orientation:   OrientationVertical,
			expectedError: false,
		},
		{
			name:          "Horizontal ship off board",
			shipType:      Battleship,
			x:             7,
			y:             0,
			orientation:   OrientationHorizontal,
			expectedError: true,
		},
		{
			name:          "Vertical ship off board",
			shipType:      Cruiser,
			x:             0,
			y:             8,
			orientation:   OrientationVertical,
			expectedError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new board for each test case
			board := NewBoard("player", "opponent")

			// Place the ship
			err := board.PlaceShip(tc.shipType, tc.x, tc.y, tc.orientation)

			// Validate the result
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify ship was added to fleet
				assert.NotNil(t, board.Fleet)
				assert.Len(t, board.Fleet, 1)

				// Verify ship details
				ship := board.Fleet[0]
				assert.Equal(t, tc.shipType, ship.ShipType)
				assert.Equal(t, tc.x, ship.Position.X)
				assert.Equal(t, tc.y, ship.Position.Y)
				assert.Equal(t, tc.orientation, ship.Orientation)
				assert.Equal(t, shipLength(tc.shipType), ship.Length)

				// Verify pins were placed correctly on the board
				length := shipLength(tc.shipType)
				for i := 0; i < length; i++ {
					checkX, checkY := tc.x, tc.y
					if tc.orientation.IsVertical() {
						checkY += i
					} else {
						checkX += i
					}
					assert.Equal(t, FieldStatePin, board.ShipsMap().FieldState(checkX, checkY),
						"Ship pin should be at (%d,%d)", checkX, checkY)
				}

				// Verify pin count was decreased
				assert.Equal(t, 30-length, board.PinsAvailable)
			}
		})
	}
}

func TestBoardPlaceOverlappingShips(t *testing.T) {
	board := NewBoard("player", "opponent")

	// Place first ship
	err := board.PlaceShip(Destroyer, 3, 3, OrientationHorizontal)
	assert.NoError(t, err)

	// Try to place second ship that overlaps with first
	err = board.PlaceShip(Cruiser, 2, 3, OrientationHorizontal)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrorIllegal)

	// Try to place second ship adjacent to first
	err = board.PlaceShip(Cruiser, 3, 2, OrientationHorizontal)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrorIllegal)
}

func TestBoardCoordinateConversion(t *testing.T) {
	board := NewBoard("player", "opponent")

	// Place a horizontal ship
	err := board.PlaceShip(Destroyer, 2, 3, OrientationHorizontal)
	assert.NoError(t, err)

	// Check that pins are placed in the correct locations
	assert.Equal(t, FieldStatePin, board.ShipsMap().FieldState(2, 3))
	assert.Equal(t, FieldStatePin, board.ShipsMap().FieldState(3, 3))
	assert.Equal(t, FieldStatePin, board.ShipsMap().FieldState(4, 3))

	// Place a vertical ship
	err = board.PlaceShip(Destroyer, 5, 5, OrientationVertical)
	assert.NoError(t, err)

	// Check that pins are placed in the correct locations
	assert.Equal(t, FieldStatePin, board.ShipsMap().FieldState(5, 5))
	assert.Equal(t, FieldStatePin, board.ShipsMap().FieldState(5, 6))
	assert.Equal(t, FieldStatePin, board.ShipsMap().FieldState(5, 7))
}

func TestBoardAttack(t *testing.T) {
	board := NewBoard("player", "opponent")

	// Place a horizontal ship
	err := board.PlaceShip(Destroyer, 2, 3, OrientationHorizontal)
	assert.NoError(t, err)

	// Attack and hit ship
	result, err := board.Attack(2, 3)
	assert.NoError(t, err)
	assert.Equal(t, FieldStateHit, result)
	assert.Equal(t, FieldStateHit, board.ShipsMap().FieldState(2, 3))

	// Attack and miss
	result, err = board.Attack(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, FieldStateMiss, result)

	// Attack the same spot again (should error)
	_, err = board.Attack(2, 3)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrorIllegal)

	// Attack off the board
	_, err = board.Attack(10, 10)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrorInvalid)
}

func TestBoardShipSinking(t *testing.T) {
	board := NewBoard("player", "opponent")

	// Place a submarine (length 2)
	err := board.PlaceShip(Submarine, 2, 3, OrientationHorizontal)
	assert.NoError(t, err)
	assert.Len(t, board.Fleet, 1)

	// Hit the ship once
	result, err := board.Attack(2, 3)
	assert.NoError(t, err)
	assert.Equal(t, FieldStateHit, result)
	
	// Check that ship has first position hit
	ship := board.Fleet[0]
	assert.True(t, ship.Hits[0])
	assert.False(t, ship.Hits[1])
	assert.False(t, ship.IsSunk())
	assert.Len(t, board.Fleet, 1) // Ship still in fleet

	// Hit the ship again to sink it
	result, err = board.Attack(3, 3)
	assert.NoError(t, err)
	assert.Equal(t, FieldStateHit, result)
	
	// Check that ship was properly hit in both positions
	// Note: The ship object is still accessible through our variable
	assert.True(t, ship.Hits[0])
	assert.True(t, ship.Hits[1])
	assert.True(t, ship.IsSunk())
	
	// With our fix to Board.Attack, the fleet should now be empty
	assert.Len(t, board.Fleet, 0, "Fleet should be empty after sinking all ships")
}

func TestBoardRemoveShip(t *testing.T) {
	board := NewBoard("player", "opponent")

	// Place a destroyer
	initialPins := board.PinsAvailable
	err := board.PlaceShip(Destroyer, 2, 3, OrientationHorizontal)
	assert.NoError(t, err)
	assert.Len(t, board.Fleet, 1)
	assert.Equal(t, initialPins-3, board.PinsAvailable)

	// Remove the ship
	err = board.RemoveShip(2, 3)
	assert.NoError(t, err)
	assert.Len(t, board.Fleet, 0)
	assert.Equal(t, initialPins, board.PinsAvailable) // Pins returned to available pool

	// Check that all ship positions are cleared
	assert.Equal(t, FieldStateEmpty, board.ShipsMap().FieldState(2, 3))
	assert.Equal(t, FieldStateEmpty, board.ShipsMap().FieldState(3, 3))
	assert.Equal(t, FieldStateEmpty, board.ShipsMap().FieldState(4, 3))

	// Try to remove a ship that doesn't exist
	err = board.RemoveShip(0, 0)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrorNotFound)
}
