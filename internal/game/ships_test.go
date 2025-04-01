package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShipsFilter(t *testing.T) {
	// Create test ships
	battleship := &Ship{ShipType: Battleship, Position: ShipPosition{X: 0, Y: 0}, Length: 5, Orientation: OrientationHorizontal}
	cruiser := &Ship{ShipType: Cruiser, Position: ShipPosition{X: 0, Y: 2}, Length: 4, Orientation: OrientationHorizontal}
	destroyer1 := &Ship{ShipType: Destroyer, Position: ShipPosition{X: 0, Y: 4}, Length: 3, Orientation: OrientationHorizontal}
	destroyer2 := &Ship{ShipType: Destroyer, Position: ShipPosition{X: 5, Y: 0}, Length: 3, Orientation: OrientationVertical}
	submarine := &Ship{ShipType: Submarine, Position: ShipPosition{X: 8, Y: 0}, Length: 2, Orientation: OrientationVertical}

	// Create fleet with all ships
	fleet := Ships{battleship, cruiser, destroyer1, destroyer2, submarine}

	// Test filtering by ship type
	battleships := fleet.Filter(byShipType(Battleship))
	assert.Len(t, battleships, 1)
	assert.Equal(t, battleship, battleships[0])

	destroyers := fleet.Filter(byShipType(Destroyer))
	assert.Len(t, destroyers, 2)
	assert.Contains(t, destroyers, destroyer1)
	assert.Contains(t, destroyers, destroyer2)

	// Test filtering by position
	shipsAtOrigin := fleet.Filter(byPosition(0, 0))
	assert.Len(t, shipsAtOrigin, 1)
	assert.Equal(t, battleship, shipsAtOrigin[0])

	// Test position within ship bounds (horizontal)
	shipsAt1_0 := fleet.Filter(byPosition(1, 0))
	assert.Len(t, shipsAt1_0, 1)
	assert.Equal(t, battleship, shipsAt1_0[0])

	// Test position within ship bounds (vertical)
	shipsAt5_1 := fleet.Filter(byPosition(5, 1))
	assert.Len(t, shipsAt5_1, 1)
	assert.Equal(t, destroyer2, shipsAt5_1[0])

	// Test position with no ships
	shipsAt9_9 := fleet.Filter(byPosition(9, 9))
	assert.Len(t, shipsAt9_9, 0)
}

func TestShipsRemove(t *testing.T) {
	// Create test ships
	battleship := &Ship{ShipType: Battleship, Position: ShipPosition{X: 0, Y: 0}, Length: 5, Orientation: OrientationHorizontal}
	cruiser := &Ship{ShipType: Cruiser, Position: ShipPosition{X: 0, Y: 2}, Length: 4, Orientation: OrientationHorizontal}
	destroyer := &Ship{ShipType: Destroyer, Position: ShipPosition{X: 0, Y: 4}, Length: 3, Orientation: OrientationHorizontal}

	// Create fleet with all ships
	initialFleet := Ships{battleship, cruiser, destroyer}
	assert.Len(t, initialFleet, 3)

	// Remove a specific ship using theShip predicate
	updatedFleet := initialFleet.Remove(theShip(cruiser))
	assert.Len(t, updatedFleet, 2)
	assert.Contains(t, updatedFleet, battleship)
	assert.Contains(t, updatedFleet, destroyer)
	assert.NotContains(t, updatedFleet, cruiser)

	// Remove by type
	updatedFleet = initialFleet.Remove(byShipType(Battleship))
	assert.Len(t, updatedFleet, 2)
	assert.NotContains(t, updatedFleet, battleship)
	assert.Contains(t, updatedFleet, cruiser)
	assert.Contains(t, updatedFleet, destroyer)

	// Remove by position
	updatedFleet = initialFleet.Remove(byPosition(0, 4))
	assert.Len(t, updatedFleet, 2)
	assert.Contains(t, updatedFleet, battleship)
	assert.Contains(t, updatedFleet, cruiser)
	assert.NotContains(t, updatedFleet, destroyer)

	// Remove all ships of a specific type
	destroyerSubmarineFleet := Ships{
		&Ship{ShipType: Destroyer, Position: ShipPosition{X: 0, Y: 0}},
		&Ship{ShipType: Destroyer, Position: ShipPosition{X: 5, Y: 0}},
		&Ship{ShipType: Submarine, Position: ShipPosition{X: 0, Y: 5}},
	}
	
	updatedFleet = destroyerSubmarineFleet.Remove(byShipType(Destroyer))
	assert.Len(t, updatedFleet, 1)
	assert.Equal(t, Submarine, updatedFleet[0].ShipType)
}

func TestShipTypeAllowances(t *testing.T) {
	// Verify the ship allowances
	assert.Equal(t, 1, shipsAllowed[Battleship])
	assert.Equal(t, 2, shipsAllowed[Cruiser])
	assert.Equal(t, 3, shipsAllowed[Destroyer])
	assert.Equal(t, 4, shipsAllowed[Submarine])
	assert.Equal(t, FleetSizeAllowed, 10) // Verify total fleet size
}

func TestShipPredicates(t *testing.T) {
	// Test byShipType predicate
	ship := &Ship{ShipType: Battleship}
	predicate := byShipType(Battleship)
	assert.True(t, predicate(ship))
	
	predicate = byShipType(Cruiser)
	assert.False(t, predicate(ship))
	
	// Test byPosition predicate
	ship = &Ship{
		ShipType:    Battleship,
		Position:    ShipPosition{X: 3, Y: 4},
		Length:      5,
		Orientation: OrientationHorizontal,
	}
	
	// Test start position
	predicate = byPosition(3, 4)
	assert.True(t, predicate(ship))
	
	// Test middle position
	predicate = byPosition(5, 4)
	assert.True(t, predicate(ship))
	
	// Test end position
	predicate = byPosition(7, 4)
	assert.True(t, predicate(ship))
	
	// Test position just outside
	predicate = byPosition(8, 4)
	assert.False(t, predicate(ship))
	
	// Test vertical ship
	ship = &Ship{
		ShipType:    Cruiser,
		Position:    ShipPosition{X: 3, Y: 4},
		Length:      4,
		Orientation: OrientationVertical,
	}
	
	// Test positions along the ship
	assert.True(t, byPosition(3, 4)(ship))
	assert.True(t, byPosition(3, 5)(ship))
	assert.True(t, byPosition(3, 6)(ship))
	assert.True(t, byPosition(3, 7)(ship))
	assert.False(t, byPosition(3, 8)(ship))
	assert.False(t, byPosition(4, 4)(ship))
	
	// Test theShip predicate
	ship1 := &Ship{ShipType: Battleship}
	ship2 := &Ship{ShipType: Cruiser}
	
	predicate = theShip(ship1)
	assert.True(t, predicate(ship1))
	assert.False(t, predicate(ship2))
}