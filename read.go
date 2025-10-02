package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"slices"
)

type Vec2 struct {
	X int
	Y int
}

type HexagonScreenInfo struct {
	Center                 Vec2
	OccupiedCheckPositions [6]Vec2 // positions to check if the side is occupied or not
}

// Define the screen positions for each hexagon in the puzzle
var HexagonScreenPositions = map[HexagonPosition]HexagonScreenInfo{
	Center: {
		Center: Vec2{X: 1310, Y: 496},
		OccupiedCheckPositions: [6]Vec2{
			{X: 1351, Y: 423}, // 0
			{X: 1391, Y: 501}, // 1
			{X: 1348, Y: 570}, // 2
			{X: 1270, Y: 562}, // 3
			{X: 1234, Y: 489}, // 4
			{X: 1273, Y: 421}, // 5
		},
	},
	TopRight: {
		Center: Vec2{X: 1435, Y: 292},
		OccupiedCheckPositions: [6]Vec2{
			{X: 1485, Y: 212}, // 0
			{X: 1527, Y: 289}, // 1
			{X: 1476, Y: 386}, // 2
			{X: 1387, Y: 367}, // 3
			{X: 1348, Y: 292}, // 4
			{X: 1394, Y: 215}, // 5
		},
	},
	Right: {
		Center: Vec2{X: 1552, Y: 511},
		OccupiedCheckPositions: [6]Vec2{
			{X: 1604, Y: 434}, // 0
			{X: 1650, Y: 518}, // 1
			{X: 1596, Y: 593}, // 2
			{X: 1500, Y: 583}, // 3
			{X: 1459, Y: 505}, // 4
			{X: 1508, Y: 431}, // 5
		},
	},
	BottomRight: {
		Center: Vec2{X: 1416, Y: 707},
		OccupiedCheckPositions: [6]Vec2{
			{X: 1462, Y: 640}, // 0
			{X: 1502, Y: 720}, // 1
			{X: 1456, Y: 786}, // 2
			{X: 1373, Y: 770}, // 3
			{X: 1335, Y: 695}, // 4
			{X: 1378, Y: 629}, // 5
		},
	},
	BottomLeft: {
		Center: Vec2{X: 1201, Y: 676},
		OccupiedCheckPositions: [6]Vec2{
			{X: 1240, Y: 614}, // 0
			{X: 1273, Y: 686}, // 1
			{X: 1235, Y: 747}, // 2
			{X: 1165, Y: 736}, // 3
			{X: 1133, Y: 667}, // 4
			{X: 1168, Y: 606}, // 5
		},
	},
	Left: {
		Center: Vec2{X: 1108, Y: 482},
		OccupiedCheckPositions: [6]Vec2{
			{X: 1142, Y: 417}, // 0
			{X: 1176, Y: 487}, // 1
			{X: 1141, Y: 551}, // 2
			{X: 1075, Y: 544}, // 3
			{X: 1043, Y: 478}, // 4
			{X: 1075, Y: 414}, // 5
		},
	},
	TopLeft: {
		Center: Vec2{X: 1210, Y: 294},
		OccupiedCheckPositions: [6]Vec2{
			{X: 1250, Y: 220}, // 0
			{X: 1286, Y: 292}, // 1
			{X: 1243, Y: 364}, // 2
			{X: 1172, Y: 365}, // 3
			{X: 1140, Y: 296}, // 4
			{X: 1177, Y: 224}, // 5
		},
	},
}

func isOccupied(col color.Color) bool {

	rgbToHsv := func(c color.Color) (h, s, v float64) {
		r, g, b, _ := c.RGBA()
		rf := float64(r) / 65535.0
		gf := float64(g) / 65535.0
		bf := float64(b) / 65535.0
		max := rf
		min := rf
		if gf > max {
			max = gf
		}
		if bf > max {
			max = bf
		}
		if gf < min {
			min = gf
		}
		if bf < min {
			min = bf
		}
		v = max
		if max == 0 {
			s = 0
		} else {
			s = (max - min) / max
		}
		if max == min {
			h = 0 // undefined
		} else if max == rf {
			h = (60*(gf-bf)/(max-min) + 360)
		} else if max == gf {
			h = (60*(bf-rf)/(max-min) + 120)
		} else {
			h = (60*(rf-gf)/(max-min) + 240)
		}
		h = h / 360.0
		return h, s, v
	}

	_, s, v := rgbToHsv(col)
	// occupied if the color is not close to white
	return (v > 0.8 && s < 0.05)

}

func ReadPuzzleFromImg(img image.Image) (*Puzzle, error) {

	NewGame := &Puzzle{}

	// debugging aid:
	// read the img and put a green pixel at the center of each hexagon
	// and a red pixel at each occupied check position if it is occupied
	// then save the img to a file
	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{0, 0}, draw.Src)
	for pos, screenInfo := range HexagonScreenPositions {
		info := HexagonScreenPositions[pos]
		dst.Set(screenInfo.Center.X, screenInfo.Center.Y, color.RGBA{0, 255, 0, 255}) // Green for center
		hexIsOccupied := isOccupied(img.At(screenInfo.Center.X, screenInfo.Center.Y))
		hex := Hexagon{
			Position: pos,
			Occupied: make([]bool, 6),
			InUse:    hexIsOccupied,
		}
		for i := range hex.Occupied {
			hex.Occupied[i] = false
		}
		NewGame.Hexagons = append(NewGame.Hexagons, hex)
		if hexIsOccupied {
			for i, checkPos := range info.OccupiedCheckPositions {
				occupied := isOccupied(img.At(checkPos.X, checkPos.Y))
				if occupied {
					dst.Set(checkPos.X, checkPos.Y, color.RGBA{255, 0, 0, 255}) // Red for occupied check positions
					NewGame.Hexagons[len(NewGame.Hexagons)-1].Occupied[i] = true
				} else {
					dst.Set(checkPos.X, checkPos.Y, color.RGBA{0, 0, 255, 255}) // Blue for unoccupied check positions
					NewGame.Hexagons[len(NewGame.Hexagons)-1].Occupied[i] = false
				}
			}
		}
	}

	// order the hexagon array by position for easier access later
	slices.SortFunc(NewGame.Hexagons, func(a, b Hexagon) int {
		return int(a.Position) - int(b.Position)
	})

	// save to file for debugging
	outFile, err := os.Create("debug_output.png")
	if err != nil {
		return nil, err
	}
	defer outFile.Close()
	png.Encode(outFile, dst)

	return NewGame, nil
}
