package main

type HexagonPosition int

const (
	None HexagonPosition = iota
	Center
	TopRight
	Right
	BottomRight
	BottomLeft
	Left
	TopLeft
)

type Hexagon struct {
	Position      HexagonPosition
	Occupied      []bool // each side can be occupied or not
	RotationCount int    // how many times it has been rotated clockwise
	InUse         bool   // whether this hexagon is part of the puzzle or not
}

// Will rotate the hexagon clockwise, changing which sides are occupied
// Warframe only allows for Clockwise rotation
func (hex *Hexagon) Rotate() {
	if !hex.InUse {
		return
	}
	last := hex.Occupied[len(hex.Occupied)-1]
	for i := len(hex.Occupied) - 1; i > 0; i-- {
		hex.Occupied[i] = hex.Occupied[i-1]
	}
	hex.Occupied[0] = last
	hex.RotationCount = (hex.RotationCount + 1) % 6
}

func (hex *Hexagon) Clone() Hexagon {
	newHex := Hexagon{
		Position:      hex.Position,
		Occupied:      make([]bool, len(hex.Occupied)),
		RotationCount: hex.RotationCount,
		InUse:         hex.InUse,
	}
	copy(newHex.Occupied, hex.Occupied)
	return newHex
}

func (hex *Hexagon) toString() string {
	s := ""
	for i, occ := range hex.Occupied {
		if occ {
			s += "1"
		} else {
			s += "0"
		}
		if i < len(hex.Occupied)-1 {
			s += ","
		}
	}
	return s
}
