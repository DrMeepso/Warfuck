package main

import (
	_ "image/png"
	"time"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func main() {

	println("Press Shift+Y to solve the puzzle on screen")

	hook.Register(hook.KeyDown, []string{"y", "shift"}, func(e hook.Event) {
		solveScreen()
	})

	s := hook.Start()
	<-hook.Process(s)

}

func solveScreen() {

	img, err := robotgo.CaptureImg()
	if err != nil {
		panic(err)
	}

	p, err := ReadPuzzleFromImg(img)
	if err != nil {
		panic(err)
	}

	solved := Solve(p)
	if solved == nil {
		println("No solution found")
		return
	}
	// print the rotation count of each hexagon in the solved puzzle
	for _, hex := range solved.Hexagons {
		if hex.InUse {
			println("Hex at position", hex.Position, "rotated", hex.RotationCount, "times")

			pos := HexagonScreenPositions[hex.Position].Center
			robotgo.Move(pos.X, pos.Y)
			time.Sleep(5 * time.Millisecond)
			for i := 0; i < hex.RotationCount; i++ {
				robotgo.Click("left")
				time.Sleep(25 * time.Millisecond)
			}

		}
	}

}
