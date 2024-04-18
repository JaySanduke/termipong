package main

import (
	"fmt"
	"time"
)

type terminalSize struct {
	width  int
	height int
}

type Coordinates struct {
	X, Y float64
}

var TerminalSize = terminalSize{
	width:  160,
	height: 40,
}

func CollisionDetection(ball *Ball, board *Board) bool {
	var ballCornerTL = ball.GetCornerPointer("TL")
	var ballCornerBR = ball.GetCornerPointer("BR")

	var ledge = ballCornerTL.X
	var redge = ballCornerBR.X
	var tedge = ballCornerTL.Y
	var bedge = ballCornerBR.Y

	if ledge <= 1 || redge >= float64(TerminalSize.width) {
		ball.Velocity.X = -ball.Velocity.X
	}

	if tedge <= 1 {
		//TODO websocket
		ball.Velocity.Y = -ball.Velocity.Y
	}

	if bedge >= float64(TerminalSize.height) {
		if ball.Position.X >= board.GetXPointer("LEFT") && ball.Position.X <= board.GetXPointer("RIGHT") {
			ball.Velocity.Y = -ball.Velocity.Y
		} else {
			ball.Velocity.Y = -ball.Velocity.Y
			//return false
			//TODO websocket - game over
		}
	}

	return true
}

func main() {
	print("\033[2J")

	var newBall = Ball{
		Height:   3,
		Width:    5,
		Position: Coordinates{X: 10, Y: 10},
		Velocity: Velocity{X: 5, Y: 3},
	}

	var newboard = Board{Height: 2, Width: 10, Position: Coordinates{X: 5, Y: 5}}

	newboard.RenderBoard(&TerminalSize)
	//newBall.RenderBall()
	go func() {
		isGameRunning := true
		for isGameRunning {
			newBall.MoveBall()
			isGameRunning = CollisionDetection(&newBall, &newboard)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	fmt.Scanln()
}
