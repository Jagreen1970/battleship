package game

import (
	"fmt"
)

type Board struct {
	PinsAvailable int
	ShotsMap      [10][10]FieldState `json:"shots_map"`
	ShipsMap      [10][10]FieldState `json:"ships_map"`
	Fleet         Ships              `json:"fleet"`
}

func NewBoard() *Board {
	return &Board{
		PinsAvailable: 30,
		ShotsMap:      [10][10]FieldState{},
		ShipsMap:      [10][10]FieldState{},
	}
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
		b.Fleet.Remove(ship(s))
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

func (b *Board) isLegalPlacement(x int, y int) bool {
	if b.isIllegalPinPosition(x, y) {
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

func (b *Board) isIllegalPinPosition(x int, y int) bool {
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
