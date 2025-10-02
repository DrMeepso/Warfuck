package main

type Puzzle struct {
	Hexagons []Hexagon
}

func (p *Puzzle) GetHexagonAtPosition(pos HexagonPosition) *Hexagon {
	return &p.Hexagons[pos-1]
}

// Center: 0
// TopRight: 1
// Right: 2
// BottomRight: 3
// BottomLeft: 4
// Left: 5
// TopLeft: 6

// get the face on the adjacent hexagon, add 3 to the index that it is returned from.

// returns the positions of adjacent hexagons in clockwise order starting from the top-right
func (p *Puzzle) GetAdjacentHexagonPositions(pos HexagonPosition) [6]HexagonPosition {
	switch pos {
	case Center:
		return [6]HexagonPosition{TopRight, Right, BottomRight, BottomLeft, Left, TopLeft}
	case TopRight:
		return [6]HexagonPosition{None, None, Right, Center, TopLeft, None}
	case Right:
		return [6]HexagonPosition{None, None, None, BottomRight, Center, TopRight}
	case BottomRight:
		return [6]HexagonPosition{Right, None, None, None, BottomLeft, Center}
	case BottomLeft:
		return [6]HexagonPosition{Center, BottomRight, None, None, None, Left}
	case Left:
		return [6]HexagonPosition{TopLeft, Center, BottomLeft, None, None, None}
	case TopLeft:
		return [6]HexagonPosition{None, TopRight, Center, Left, None, None}
	default:
		return [6]HexagonPosition{None, None, None, None, None, None}
	}
}

func (p *Puzzle) IsSolved() bool {

	for i := 0; i < len(p.Hexagons); i++ {
		hex := &p.Hexagons[i]
		if !hex.InUse {
			continue
		}
		//println("Checking hex at position", hex.Position)
		//println(hex.toString())
		adjacentPositions := p.GetAdjacentHexagonPositions(hex.Position)
		for i, adjPos := range adjacentPositions {
			// print the raw loop index to see order 0..5
			//println(" Checking: ", hex.Position, "side", i, " on adjacent hex at position", pos, "which is side", (i+3)%6)
			var hexAtPos *Hexagon
			if adjPos != None {
				hexAtPos = p.GetHexagonAtPosition(adjPos)
				//println("	", (i+3)%6, "vs", i, ":", hexAtPos.Occupied[(i+3)%6], "vs", hex.Occupied[i], "at", hex.Position, "and", hexAtPos.Position)
				//println("  ", hexAtPos.toString(), "vs", hex.toString())
			}
			if adjPos != None && !hexAtPos.Occupied[(i+3)%6] && hex.Occupied[i] {
				//println("Failed at hex", hex.Position, "side", i, "adjacent hex", adjPos, "side", (i+3)%6)
				return false
			}
			if adjPos == None && hex.Occupied[i] {
				//println("Failed at hex", hex.Position, "side", i, "which has no adjacent hexagon")
				return false
			}
		}
	}
	return true
}

func (p *Puzzle) Clone() *Puzzle {
	newPuzzle := &Puzzle{
		Hexagons: make([]Hexagon, len(p.Hexagons)),
	}
	for i, hex := range p.Hexagons {
		newPuzzle.Hexagons[i] = hex.Clone()
	}
	return newPuzzle
}
