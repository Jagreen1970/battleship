package battleship

import (
	"fmt"
	"math"
	"sort"
)

type ShipPart struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	State FieldState
}

func NewShipPart(x int, y int, state FieldState) *ShipPart {
	return &ShipPart{X: x, Y: y, State: state}
}

func (p *ShipPart) Is(x int, y int) bool {
	return p.X == x && p.Y == y
}

func (p *ShipPart) IsNextTo(x int, y int) bool {
	return math.Abs(float64(p.X-x)) == 1 || math.Abs(float64(p.Y-y)) == 1
}

type ShipParts []*ShipPart

func (pts ShipParts) All(predicate func(part *ShipPart) bool) bool {
	for _, part := range pts {
		if !predicate(part) {
			return false
		}
	}
	return true
}

func stateIs(state FieldState) func(part *ShipPart) bool {
	return func(part *ShipPart) bool {
		return part.State == state
	}
}

type ShipType string

type ShipOrientation string

const (
	BattleShip            ShipType        = "Battleship"
	Cruiser               ShipType        = "Cruiser"
	Destroyer             ShipType        = "Destroyer"
	Submarine             ShipType        = "Submarine"
	UnknownShip           ShipType        = "Unknown"
	InvalidShip           ShipType        = "Invalid"
	OrientationHorizontal ShipOrientation = "Horizontal"
	OrientationVertical   ShipOrientation = "Vertical"
	OrientationUnknown    ShipOrientation = "Unknown"
	OrientationInvalid    ShipOrientation = "Invalid"
)

type Ship struct {
	ShipType    ShipType        `json:"ship_type"`
	Parts       ShipParts       `json:"fields"`
	Orientation ShipOrientation `json:"orientation"`
}

func NewShip(x int, y int) *Ship {
	return &Ship{
		ShipType: UnknownShip,
		Parts: ShipParts{
			NewShipPart(x, y, FieldStatePin),
		},
	}
}

func NewShipWithParts(parts ShipParts) *Ship {
	return &Ship{
		ShipType: UnknownShip,
		Parts:    parts,
	}
}

func (s *Ship) Hit(x int, y int) bool {
	for _, part := range s.Parts {
		if part.Is(x, y) {
			s.takeHit(part)
		}
	}
	return s.IsSunk()
}

func (s *Ship) IsAtPosition(x int, y int) bool {
	for _, field := range s.Parts {
		if field.Is(x, y) {
			return true
		}
	}
	return false
}

func (s *Ship) IsNextToPosition(x int, y int) bool {
	for _, field := range s.Parts {
		if field.IsNextTo(x, y) {
			return true
		}
	}
	return false
}

func (s *Ship) CanMergeTo(ships ...*Ship) bool {
	mergedShip := Ship{
		Parts: s.Parts,
	}
	for _, ship := range ships {
		mergedShip.Parts = append(mergedShip.Parts, ship.Parts...)
	}

	mergedShip.AdjustProperties()

	return mergedShip.IsValid()
}

func (s *Ship) Merge(ships ...*Ship) (*Ship, error) {
	mergedShip := &Ship{
		Parts: s.Parts,
	}
	for _, ship := range ships {
		mergedShip.Parts = append(mergedShip.Parts, ship.Parts...)
	}

	mergedShip.AdjustProperties()
	if !mergedShip.IsValid() {
		return s, fmt.Errorf("cannot merge %v to %v: %w", s, ships, ErrorIllegal)
	}

	return mergedShip, nil
}

func (s *Ship) AdjustProperties() {
	s.sortFields()
	s.adjustOrientation()
	s.adjustShipType()
}

func (s *Ship) IsValid() bool {
	if s.ShipType == UnknownShip || s.ShipType == InvalidShip {
		return false
	}

	if s.Orientation == OrientationInvalid || s.Orientation == OrientationUnknown {
		return false
	}

	return s.hasNoGaps()
}

func (s *Ship) IsVertical() bool {
	for i := 1; i < len(s.Parts); i++ {
		if s.Parts[i].X != s.Parts[i-1].X {
			return false
		}
	}
	return len(s.Parts) > 1
}

func (s *Ship) IsHorizontal() bool {
	for i := 1; i < len(s.Parts); i++ {
		if s.Parts[i].Y != s.Parts[i-1].Y {
			return false
		}
	}
	return len(s.Parts) > 1
}

func (s *Ship) IsSunk() bool {
	return s.Parts.All(stateIs(FieldStateHit))
}

func (s *Ship) adjustShipType() {
	var newType ShipType
	switch len(s.Parts) {
	case 5:
		newType = BattleShip
	case 4:
		newType = Cruiser
	case 3:
		newType = Destroyer
	case 2:
		newType = Submarine
	case 1:
		newType = UnknownShip
	default:
		newType = InvalidShip
	}

	s.ShipType = newType
}

func (s *Ship) adjustOrientation() {
	if len(s.Parts) < 1 {
		s.Orientation = OrientationInvalid
		return
	}

	if s.IsVertical() {
		s.Orientation = OrientationVertical
		return
	}
	if s.IsHorizontal() {
		s.Orientation = OrientationHorizontal
		return
	}

	s.Orientation = OrientationUnknown
}

func (s *Ship) sortFields() {
	sort.Slice(s.Parts, func(i, j int) bool {
		if s.Parts[i].Y < s.Parts[j].Y {
			return true
		}
		if s.Parts[i].Y > s.Parts[j].Y {
			return false
		}
		return s.Parts[i].X < s.Parts[j].X
	})
}

func (s *Ship) hasNoGaps() bool {
	if s.Orientation == OrientationHorizontal {
		return s.hasNoHorizontalGap()
	}
	if s.Orientation == OrientationVertical {
		return s.hasNoVerticalGap()
	}
	return false
}

func (s *Ship) hasNoHorizontalGap() bool {
	for i := 1; i < len(s.Parts); i++ {
		if s.Parts[i-1].X != s.Parts[i].X-1 {
			return false
		}
	}
	return true
}

func (s *Ship) hasNoVerticalGap() bool {
	for i := 1; i < len(s.Parts); i++ {
		if s.Parts[i-1].Y != s.Parts[i].Y-1 {
			return false
		}
	}
	return true
}

func (s *Ship) takeHit(shipPart *ShipPart) {
	shipPart.State = FieldStateHit
}
