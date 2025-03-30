package game

type ShipPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p ShipPosition) Is(x int, y int) bool {
	return p.X == x && p.Y == y
}

func (p ShipPosition) IsNextTo(x int, y int) bool {
	return (p.X == x && (p.Y == y+1 || p.Y == y-1)) || // vertically adjacent
		(p.Y == y && (p.X == x+1 || p.X == x-1)) // horizontally adjacent
}

type ShipType string

type ShipOrientation string

func (o ShipOrientation) IsVertical() bool {
	return o == OrientationVertical
}

const (
	Battleship            ShipType        = "Battleship"
	Cruiser               ShipType        = "Cruiser"
	Destroyer             ShipType        = "Destroyer"
	Submarine             ShipType        = "Submarine"
	InvalidShip           ShipType        = "Invalid"
	OrientationHorizontal ShipOrientation = "Horizontal"
	OrientationVertical   ShipOrientation = "Vertical"
)

type Ship struct {
	ShipType    ShipType        `json:"ship_type"`
	Position    ShipPosition    `json:"position"`
	Length      int             `json:"length"`
	Hits        []bool          `json:"hits"`
	Orientation ShipOrientation `json:"orientation"`
}

func NewShip(shipType ShipType, x int, y int, isVertical bool) *Ship {
	length := shipLength(shipType)
	orientation := OrientationHorizontal
	if isVertical {
		orientation = OrientationVertical
	}

	return &Ship{
		ShipType:    shipType,
		Position:    ShipPosition{X: x, Y: y},
		Length:      length,
		Hits:        make([]bool, length),
		Orientation: orientation,
	}
}

func (s *Ship) Hit(x int, y int) bool {
	if !s.IsAtPosition(x, y) {
		return false
	}

	index := s.getPartIndex(x, y)
	if index >= 0 {
		s.Hits[index] = true
	}
	return s.IsSunk()
}

func (s *Ship) getPartIndex(x, y int) int {
	if s.Orientation.IsVertical() {
		return y - s.Position.Y
	}
	return x - s.Position.X
}

func (s *Ship) IsAtPosition(x int, y int) bool {
	if s.Orientation.IsVertical() {
		return x == s.Position.X && y >= s.Position.Y && y < s.Position.Y+s.Length
	}
	return y == s.Position.Y && x >= s.Position.X && x < s.Position.X+s.Length
}

func (s *Ship) IsNextToPosition(x int, y int) bool {
	// Check each position of the ship
	for i := 0; i < s.Length; i++ {
		checkX := s.Position.X
		checkY := s.Position.Y
		if s.Orientation.IsVertical() {
			checkY += i
		} else {
			checkX += i
		}

		pos := ShipPosition{X: checkX, Y: checkY}
		if pos.IsNextTo(x, y) {
			return true
		}
	}
	return false
}

func (s *Ship) IsSunk() bool {
	for _, hit := range s.Hits {
		if !hit {
			return false
		}
	}
	return true
}

func shipLength(shipType ShipType) int {
	switch shipType {
	case Battleship:
		return 5
	case Cruiser:
		return 4
	case Destroyer:
		return 3
	case Submarine:
		return 2
	default:
		return 0
	}
}
