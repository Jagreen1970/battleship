package battleship

import (
	"errors"
	"fmt"
	"log"
)

type Board struct {
	PinsAvailable int                `json:"pins_available" bson:"pins_available"`
	ShotsMap      [10][10]FieldState `json:"shots_map" bson:"shots_map"`
	ShipsMap      [10][10]FieldState `json:"ships_map" bson:"ships_map"`
	Fleet         Ships              `json:"fleet" bson:"fleet"`
}

func NewBoard() *Board {
	return &Board{
		PinsAvailable: 30,
		ShotsMap:      [10][10]FieldState{},
		ShipsMap:      [10][10]FieldState{},
	}
}

// ValidSetup checks if the board is in a valid setup state. Not all pins need to be placed yet, but any placed pin has to be
// in a valid position.
func (b *Board) ValidSetup() error {
	if len(b.Fleet) == 0 {
		return nil
	}

	if len(b.Fleet) > FleetSizeAllowed {
		return fmt.Errorf("too many ships: %d (%w)", len(b.Fleet), ErrorIllegal)
	}

	for shipType, numAllowed := range shipsAllowed {
		if shipType != UnknownShip && len(b.Fleet.Filter(byShipType(shipType))) > numAllowed {
			return fmt.Errorf("too many ships of type: %q (%w)", shipType, ErrorIllegal)
		}
	}

	return nil
}

func (b *Board) CanAttack(x int, y int) error {
	if b.offBoard(x, y) {
		return fmt.Errorf("shot (%d, %d) is off board: %w", x, y, ErrorInvalid)
	}

	if b.alreadyTried(x, y) {
		return fmt.Errorf("already tried to shoot at (%d, %d): %w", x, y, ErrorIllegal)
	}
	return nil
}

func (b *Board) Attack(x int, y int) (FieldState, error) {
	if b.ShipsMap[x][y] != FieldStatePin {
		return FielStateMiss, nil
	}

	s, err := b.ShipAtPosition(x, y)
	if err != nil {
		return FieldStateUnknown, fmt.Errorf("no ship found at position (%d, %d): %w", x, y, err)
	}

	sunk := s.Hit(x, y)
	if sunk {
		b.Fleet.Remove(theShip(s))
	}
	b.ShipsMap[x][y] = FieldStateHit
	return FieldStateHit, nil
}

func (b *Board) Track(fieldState FieldState, x int, y int) {
	b.ShotsMap[x][y] = fieldState
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

func (b *Board) PlacePin(x int, y int) error {
	if b.PinsAvailable <= 0 {
		return fmt.Errorf("you have already used all your pins (%w)", ErrorInvalid)
	}

	if !b.isLegalPlacement(x, y) {
		return fmt.Errorf("you are not allowed to place a pin in position: %d, %d. (%w)", x, y, ErrorIllegal)
	}

	b.PinsAvailable--
	b.ShipsMap[x][y] = FieldStatePin

	if b.isolatedPlacement(x, y) {
		b.addNewShip(x, y)
		return nil
	}

	return b.mergeToNewShip(x, y)
}

func (b *Board) RecoverPin(x int, y int) error {
	if b.PinsAvailable >= 30 {
		return fmt.Errorf("you have already recovered all your pins (%w)", ErrorInvalid)
	}

	if !b.isLegalRecovery(x, y) {
		return fmt.Errorf("you are not allowed to recover the pin in position: %d, %d. (%w)", x, y, ErrorIllegal)
	}
	b.PinsAvailable++
	b.ShipsMap[x][y] = FieldStateEmpty

	if b.isolatedPlacement(x, y) {
		b.removeShip(x, y)
		return nil
	}

	return b.shortenOrSplitShip(x, y)
}

func (b *Board) isLegalRecovery(x int, y int) bool {
	if b.ShipsMap[x][y] != FieldStatePin {
		return false
	}

	if b.isolatedPlacement(x, y) {
		return true
	}

	return b.canShortenOrSplitShip(x, y)
}

func (b *Board) canShortenOrSplitShip(x int, y int) bool {
	ships := b.Fleet.Filter(byPosition(x, y))
	l := len(ships)
	if l != 1 {
		return false
	}

	s := ships[0]
	if s.ShipType == UnknownShip || s.ShipType == InvalidShip {
		return false
	}

	return true
}

func (b *Board) shortenOrSplitShip(x int, y int) error {
	ships := b.Fleet.Filter(byPosition(x, y))
	l := len(ships)
	if l != 1 {
		return fmt.Errorf("cannot shorten or split ship: must have exactly 1 ship")
	}

	s := ships[0]
	if s.ShipType == InvalidShip || s.ShipType == UnknownShip {
		return fmt.Errorf("cannot shorten or split ship: invalid ship type")
	}

	b.removeShip(x, y)
	for i, part := range s.Parts {
		if !part.Is(x, y) {
			continue
		}

		p1, p2 := s.Parts[:i+1], s.Parts[i+1:]
		if len(p1) != 0 {
			s1 := NewShipWithParts(p1)
			s1.AdjustProperties()
			b.Fleet = append(b.Fleet, s1)
		}
		if len(p2) != 0 {
			s2 := NewShipWithParts(p2)
			s2.AdjustProperties()
			b.Fleet = append(b.Fleet, s2)
		}
		break
	}

	return nil
}

func (b *Board) isLegalPlacement(x int, y int) bool {
	if b.isIllegalPinPlacement(x, y) {
		return false
	}

	if b.isolatedPlacement(x, y) {
		if b.canAddShip() {
			return true
		}
		return false
	}

	return b.canMergeToNewShip(x, y)
}

func (b *Board) isIllegalPinPlacement(x int, y int) bool {
	// Pin is off-board
	if b.offBoard(x, y) {
		return true
	}

	// Position is already taken
	if b.ShipsMap[x][y] != FieldStateEmpty {
		return true
	}

	// diagonally adjacent cells must be empty
	if b.anyDiagonalsOccupied(x, y) {
		return true
	}

	return false
}

func (b *Board) isolatedPlacement(x int, y int) bool {
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if (i == x && j == y) || b.offBoard(i, j) {
				continue
			}
			if b.ShipsMap[i][j] == FieldStatePin {
				return false
			}
		}
	}
	return true
}

func (b *Board) anyDiagonalsOccupied(x int, y int) bool {
	l := x - 1
	r := x + 1
	u := y - 1
	d := y + 1

	if l >= 0 {
		if u >= 0 {
			// upper left
			if b.ShipsMap[l][u] == FieldStatePin {
				return true
			}
		}
		if d < 10 {
			// lower left
			if b.ShipsMap[l][d] == FieldStatePin {
				return true
			}
		}
	}

	if r < 10 {
		if u >= 0 {
			// upper right
			if b.ShipsMap[r][u] == FieldStatePin {
				return true
			}
		}
		if d < 10 {
			// lower right
			if b.ShipsMap[r][d] == FieldStatePin {
				return true
			}
		}
	}

	return false
}

func (b *Board) addNewShip(x int, y int) {
	b.Fleet = append(b.Fleet, NewShip(x, y))
}

func (b *Board) removeShip(x int, y int) {
	s := b.Fleet.Filter(byPosition(x, y))
	if len(s) != 1 {
		log.Fatal(errors.New("must filter exactly one ship"))
		return
	}
	b.Fleet.Remove(theShip(s[0]))
}

func (b *Board) canMergeToNewShip(x int, y int) bool {
	ships := b.Fleet.Filter(byAdjacentPosition(x, y))
	l := len(ships)
	if l > 2 || l < 1 {
		return false
	}

	return NewShip(x, y).CanMergeTo(ships...)
}

func (b *Board) canAddShip() bool {
	return len(b.Fleet) < 9
}

func (b *Board) mergeToNewShip(x int, y int) error {
	ships := b.Fleet.Filter(byAdjacentPosition(x, y))
	l := len(ships)
	if l > 2 || l < 1 {
		return fmt.Errorf("cannot merge to new ship: %w", ErrorIllegal)
	}

	newShip, err := NewShip(x, y).Merge(ships...)

	if err != nil {
		return fmt.Errorf("cannot merge to new ship: %w", err)
	}

	b.Fleet = b.Fleet.Remove(byAdjacentPosition(x, y))
	b.Fleet = append(b.Fleet, newShip)

	return nil
}

func (b *Board) offBoard(x int, y int) bool {
	return x < 0 || y < 0 || x >= 10 || y >= 10
}

func (b *Board) alreadyTried(x int, y int) bool {
	return b.ShotsMap[x][y] != FieldStateEmpty
}

func (b *Board) Lost() bool {
	return len(b.Fleet) == 0
}
