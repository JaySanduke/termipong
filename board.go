package main

import (
	"fmt"
	"math"
)

type Board struct {
	Height, Width int
	Position      Coordinates
	Speed         int
}

func (boardstate *Board) GetXPointer(direction string) float64 {
	if direction == "LEFT" {
		return boardstate.Position.X - math.Floor(float64(boardstate.Width/2))
	} else if direction == "RIGHT" {
		return boardstate.Position.X + math.Floor(float64(boardstate.Width/2))
	} else {
		return boardstate.Position.X
	}
}

func (boardstate *Board) RenderBoard(tsize *terminalSize) {
	fmt.Print("\033[" + fmt.Sprint(tsize.height+1) + ";0H")
	for i := float64(0); int(i) < boardstate.Height; i++ {
		for j := float64(0); int(j) < tsize.width; j++ {
			if j >= boardstate.GetXPointer("LEFT") && j <= boardstate.GetXPointer("RIGHT") {

				if i == float64(0) {
					if j == boardstate.GetXPointer("RIGHT") {
						fmt.Print("\\")
						continue
					} else if j == boardstate.GetXPointer("LEFT") {
						fmt.Print("/")
						continue
					}
				} else if int(i) == boardstate.Height-1 {
					if j == boardstate.GetXPointer("RIGHT") {
						fmt.Print("/")
						continue
					} else if j == boardstate.GetXPointer("LEFT") {
						fmt.Print("\\")
						continue

					}
				}
				fmt.Print("=")
				continue
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

func (boardstate *Board) MoveBoard(direction string) {

	if boardstate.Speed == 0 {
		boardstate.Speed = 1
	}

	if direction == "LEFT" {
		boardstate.Position.X -= float64(boardstate.Speed)
	} else if direction == "RIGHT" {
		boardstate.Position.X += float64(boardstate.Speed)
	}
}
