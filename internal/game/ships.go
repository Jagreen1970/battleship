package game

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

func ship(s *Ship) func(*Ship) bool {
	return func(sh *Ship) bool {
		return s == sh
	}
}

func byAdjacentPosition(x, y int) func(s *Ship) bool {
	return func(s *Ship) bool {
		return s.IsNextToPosition(x, y)
	}
}

func byPosition(x, y int) func(s *Ship) bool {
	return func(s *Ship) bool {
		return s.IsAtPosition(x, y)
	}
}
