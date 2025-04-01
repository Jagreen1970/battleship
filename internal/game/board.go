package game

import (
	"fmt"
)

type BoardMap struct {
	Title string       `json:"title"`
	Map   [10]FieldRow `json:"map"`
}

func (m *BoardMap) FieldState(x int, y int) FieldState {
	// IMPORTANT: In this codebase, the convention is that:
	// - In the Map data structure, the first index is the row (y) and the second is the column (x)
	// - But in the API and user-facing commands, we use (x,y) where x is horizontal and y is vertical
	// So we need to swap the coordinates when accessing the internal map
	return m.Map[y][x]
}

func (m *BoardMap) Set(x int, y int, fieldState FieldState) {
	// Same coordinate conversion as in FieldState
	m.Map[y][x] = fieldState
}

func (m *BoardMap) Print() {
	fmt.Println(m.Title)
	for _, row := range m.Map {
		for _, field := range row {
			fmt.Printf("%c", field)
		}
		fmt.Println()
	}
}

type Board struct {
	PinsAvailable int          `json:"pins_available" bson:"pins_available"`
	Maps          [2]*BoardMap `json:"maps" bson:"maps"`
	Fleet         Ships        `json:"fleet" bson:"fleet"`
}

func NewBoard(playerName, opponentName string) *Board {
	b := Board{
		PinsAvailable: 30,
		Maps:          [2]*BoardMap{},
		Fleet:         nil,
	}
	b.Maps[0] = &BoardMap{
		Title: playerName,
		Map:   [10]FieldRow{},
	}
	b.Maps[1] = &BoardMap{
		Title: opponentName,
		Map:   [10]FieldRow{},
	}

	for x := range 10 {
		for y := range 10 {
			b.Maps[0].Set(x, y, FieldStateEmpty)
			b.Maps[1].Set(x, y, FieldStateEmpty)
		}
	}
	return &b
}

// ValidSetup checks if the board is in a valid setup state. Not all ships need to be placed yet, but any
// placed ship has to be in a valid position.
func (b *Board) ValidSetup() error {
	if len(b.Fleet) == 0 {
		return nil
	}

	if len(b.Fleet) > FleetSizeAllowed {
		return fmt.Errorf("too many ships: %d (%w)", len(b.Fleet), ErrorIllegal)
	}

	for shipType, numAllowed := range shipsAllowed {
		if len(b.Fleet.Filter(byShipType(shipType))) > numAllowed {
			return fmt.Errorf("too many ships of type: %q (%w)", shipType, ErrorIllegal)
		}
	}

	return nil
}

// CanAttack checks if an attack can be made at the specified coordinates.
// Returns nil if the attack is valid, otherwise returns an error with the reason:
// - ErrorInvalid if coordinates are outside the 10x10 board
// - ErrorIllegal if the position was already attacked
func (b *Board) CanAttack(x int, y int) error {
	if b.offBoard(x, y) {
		return fmt.Errorf("shot (%d, %d) is off board: %w", x, y, ErrorInvalid)
	}

	if b.alreadyTried(x, y) {
		return fmt.Errorf("already tried to shoot at (%d, %d): %w", x, y, ErrorIllegal)
	}
	return nil
}

// Attack processes an attack at the specified coordinates (x,y) and returns the result.
// Returns:
// - FieldStateMiss and nil if no ship was hit
// - FieldStateHit and nil if a ship was hit
// - FieldStateUnknown and error if ship lookup fails
//
// The method:
// 1. Checks if there's a ship pin at the coordinates
// 2. Locates the ship at those coordinates
// 3. Records the hit on the ship
// 4. Removes the ship from fleet if sunk
// 5. Updates the board state
func (b *Board) Attack(x int, y int) (FieldState, error) {
	// Check if position is valid
	if b.offBoard(x, y) {
		return FieldStateUnknown, fmt.Errorf("shot (%d, %d) is off board: %w", x, y, ErrorInvalid)
	}

	// Check if the position was already hit before
	if b.ShipsMap().FieldState(x, y) == FieldStateHit {
		return FieldStateUnknown, fmt.Errorf("already tried to shoot at (%d, %d): %w", x, y, ErrorIllegal)
	}

	// Check if there's a ship pin at the coordinates
	if b.ShipsMap().FieldState(x, y) != FieldStatePin {
		return FieldStateMiss, nil
	}

	// Locate the ship at those coordinates
	s, err := b.ShipAtPosition(x, y)
	if err != nil {
		return FieldStateUnknown, fmt.Errorf("no ship found at position (%d, %d): %w", x, y, err)
	}

	sunk := s.Hit(x, y)
	if sunk {
		b.Fleet.Remove(theShip(s))
	}
	b.ShipsMap().Set(x, y, FieldStateHit)
	return FieldStateHit, nil
}

// Track records the result of an attack on the shots map at the specified coordinates.
// The fieldState parameter indicates the outcome (hit, miss, etc.) of the attack.
func (b *Board) Track(fieldState FieldState, x int, y int) {
	b.ShotsMap().Set(x, y, fieldState)
}

func (b *Board) ShipAtPosition(x int, y int) (*Ship, error) {
	ships := b.Fleet.Filter(byPosition(x, y))
	if len(ships) == 0 {
		return nil, fmt.Errorf("cannot find ship at position (%d, %d): %w", x, y, ErrorNotFound)
	}
	if len(ships) > 1 {
		return nil, fmt.Errorf("cannot find ship at position (%d, %d): %w", x, y, ErrorAmbiguous)
	}

	return ships[0], nil
}

type ShipPlacement struct {
	StartX, StartY int
	Length         int
	Orientation    ShipOrientation
}

func (b *Board) PlaceShip(shipType ShipType, x, y int, orientation ShipOrientation) error {
	length := shipLength(shipType)
	placement := ShipPlacement{
		StartX:      x,
		StartY:      y,
		Length:      length,
		Orientation: orientation,
	}

	if err := b.validateShipPlacement(placement); err != nil {
		return fmt.Errorf("invalid ship placement: %w", err)
	}

	if !b.hasAvailableShipSlot(shipType) {
		return fmt.Errorf("no more ships of type %v allowed: %w", shipType, ErrorIllegal)
	}

	ship := b.createShip(shipType, placement)
	b.Fleet = append(b.Fleet, ship)
	b.PinsAvailable -= length

	// Mark all ship positions on the board
	for i := 0; i < length; i++ {
		b.ShipsMap().Set(x, y, FieldStatePin)
		x, y = increment(x, y, orientation.IsVertical())
	}

	return nil
}

func increment(x int, y int, vertical bool) (int, int) {
	// The original issue is that when vertical=true, we should increment y,
	// and when vertical=false (horizontal), we should increment x.
	// This matches the ShipOrientation constants in ship.go.
	if vertical {
		y++
	} else {
		x++
	}
	return x, y
}

func (b *Board) clearShipFromMap(ship *Ship) {
	for i := 0; i < ship.Length; i++ {
		posX := ship.Position.X
		posY := ship.Position.Y
		if ship.Orientation.IsVertical() {
			posY += i
		} else {
			posX += i
		}
		b.ShipsMap().Set(posX, posY, FieldStateEmpty)
	}
}

func (b *Board) RemoveShip(x int, y int) error {
	ship, err := b.ShipAtPosition(x, y)
	if err != nil {
		return fmt.Errorf("cannot determine ship at position (%d, %d): %w", x, y, err)
	}

	b.clearShipFromMap(ship)

	// Remove the ship from the fleet
	b.Fleet = b.Fleet.Remove(theShip(ship))

	// Return pins to available pool
	b.PinsAvailable += ship.Length

	return nil
}

func (b *Board) validateShipPlacement(p ShipPlacement) error {
	// Check if ship would be off board
	endX := p.StartX
	endY := p.StartY
	if p.Orientation.IsVertical() {
		endY += p.Length - 1
	} else {
		endX += p.Length - 1
	}

	if b.offBoard(p.StartX, p.StartY) || b.offBoard(endX, endY) {
		return fmt.Errorf("ship placement out of bounds: %w", ErrorInvalid)
	}

	// Check if ship would overlap with other ships or their surrounding area
	for i := -1; i <= p.Length; i++ {
		for j := -1; j <= 1; j++ {
			checkX := p.StartX
			checkY := p.StartY
			if p.Orientation.IsVertical() {
				checkY += i
				checkX += j
			} else {
				checkX += i
				checkY += j
			}

			if !b.offBoard(checkX, checkY) && b.ShipsMap().FieldState(checkX, checkY) == FieldStatePin {
				return fmt.Errorf("ship would overlap with existing ship or adjacent area: %w", ErrorIllegal)
			}
		}
	}

	return nil
}

func (b *Board) hasAvailableShipSlot(shipType ShipType) bool {
	currentCount := len(b.Fleet.Filter(byShipType(shipType)))
	allowedCount, exists := shipsAllowed[shipType]
	return exists && currentCount < allowedCount
}

func (b *Board) createShip(shipType ShipType, p ShipPlacement) *Ship {
	ship := &Ship{
		ShipType:    shipType,
		Position:    ShipPosition{X: p.StartX, Y: p.StartY},
		Length:      p.Length,
		Hits:        make([]bool, p.Length),
		Orientation: p.Orientation,
	}
	return ship
}

func (b *Board) offBoard(x int, y int) bool {
	return x < 0 || y < 0 || x >= 10 || y >= 10
}

func (b *Board) alreadyTried(x int, y int) bool {
	return b.ShotsMap().FieldState(x, y) != FieldStateEmpty
}

func (b *Board) Lost() bool {
	return len(b.Fleet) == 0
}

func (b *Board) ShotsMap() *BoardMap {
	return b.Maps[1]
}

func (b *Board) ShipsMap() *BoardMap {
	return b.Maps[0]
}

func (b *Board) Print() {
	fmt.Println("Ships map:")
	b.ShipsMap().Print()
	fmt.Println("Shots map:")
	b.ShotsMap().Print()
}
