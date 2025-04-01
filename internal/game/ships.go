package game

var shipsAllowed = map[ShipType]int{
	Battleship: 1,
	Cruiser:    2,
	Destroyer:  3,
	Submarine:  4,
}

const FleetSizeAllowed = 10

type Ships []*Ship

func (f Ships) Filter(predicate func(*Ship) bool) Ships {
	var result Ships
	for _, ship := range f {
		if predicate(ship) {
			result = append(result, ship)
		}
	}
	return result
}

func (f Ships) Remove(predicate func(*Ship) bool) Ships {
	var result Ships
	for _, ship := range f {
		if !predicate(ship) {
			result = append(result, ship)
		}
	}
	return result
}

func theShip(s *Ship) func(*Ship) bool {
	return func(sh *Ship) bool {
		return s == sh
	}
}

func byPosition(x, y int) func(s *Ship) bool {
	return func(s *Ship) bool {
		return s.IsAtPosition(x, y)
	}
}

func byShipType(shipType ShipType) func(*Ship) bool {
	return func(s *Ship) bool {
		return s.ShipType == shipType
	}
}
