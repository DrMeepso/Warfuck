package main

import (
	"fmt"
	"math"
	"sync"
)

type RotateState [7]int // how many times each hexagon has been rotated
// this starts at 0 for the base state.

func Solve(p *Puzzle) *Puzzle {

	inUseHexagons := 0
	for _, hex := range p.Hexagons {
		if hex.InUse {
			inUseHexagons++
		}
	}
	totalStates := math.Pow(6, float64(inUseHexagons))
	println("Total states to check:", int(totalStates))

	EveryState := make([]RotateState, 0)

	// add all states to EveryState
	var generateStates func(currentState *RotateState, hexIndex int)
	generateStates = func(currentState *RotateState, hexIndex int) {
		if hexIndex >= len(p.Hexagons) {
			EveryState = append(EveryState, *currentState)
			return
		}
		hex := p.Hexagons[hexIndex]
		if hex.InUse {
			for i := 0; i < 6; i++ {
				currentState[hexIndex] = i
				generateStates(currentState, hexIndex+1)
			}
		} else {
			currentState[hexIndex] = 0
			generateStates(currentState, hexIndex+1)
		}
	}
	generateStates(&RotateState{0, 0, 0, 0, 0, 0, 0}, 0)

	println("Generated all states, now checking...", len(EveryState))

	var wg sync.WaitGroup
	var solvedPuzzle *Puzzle
	var mu sync.Mutex
	found := false

	for _, state := range EveryState {

		wg.Add(1)

		go func() {

			// apply the state to a clone of the puzzle
			testPuzzle := p.Clone()
			for hexIndex, rotateCount := range state {
				for r := 0; r < rotateCount; r++ {
					testPuzzle.Hexagons[hexIndex].Rotate()
				}
				testPuzzle.Hexagons[hexIndex].RotationCount = rotateCount
			}
			//println("Checking state", i, "of", len(EveryState), ":", RotationToString(state))
			if testPuzzle.IsSolved() {
				//println("Solved after checking", i, "states")
				mu.Lock()
				if !found {
					solvedPuzzle = testPuzzle
					found = true
				}
				mu.Unlock()
			}

			wg.Done()

		}()
	}

	wg.Wait()
	if solvedPuzzle != nil {
		return solvedPuzzle
	}

	println("No solution found after checking all", len(EveryState), "states")

	return nil
}

func RotationToString(r [7]int) string {
	s := ""
	println(len(r))
	for i := 0; i < len(r); i++ {
		v := r[i]
		s += "Hex " + fmt.Sprint(i) + ": " + fmt.Sprint(v) + "\n"
	}
	return s
}
