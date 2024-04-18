package main

import (
	"fmt"
	"math"
)

type Velocity struct {
	X, Y int
}

type Ball struct {
	Height, Width float64
	Position      Coordinates
	Velocity      Velocity
}

func (v *Velocity) GetUnitVector() (float64, float64) {
	magnitude := math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
	return float64(v.X) / magnitude, float64(v.Y) / magnitude
}

func (ballstate *Ball) GetCornerPointer(corner string) Coordinates {
	switch corner {
	case "TL":
		return Coordinates{X: ballstate.Position.X - ballstate.Width/2, Y: ballstate.Position.Y - ballstate.Height/2}
	case "TR":
		return Coordinates{X: ballstate.Position.X + ballstate.Width/2, Y: ballstate.Position.Y - ballstate.Height/2}
	case "BL":
		return Coordinates{X: ballstate.Position.X - ballstate.Width/2, Y: ballstate.Position.Y + ballstate.Height/2}
	case "BR":
		return Coordinates{X: ballstate.Position.X + ballstate.Width/2, Y: ballstate.Position.Y + ballstate.Height/2}
	default:
		return ballstate.Position

	}
}

func (ballstate *Ball) ClearBall() {
	var ballCorner = ballstate.GetCornerPointer("TL")

	for i := float64(0); i < ballstate.Height; i++ {
		//fmt.Print("\033[" + fmt.Sprint(ballCorner.Y+i) + ";" + fmt.Sprint(ballCorner.X) + "H")
		for j := float64(0); j < ballstate.Width; j++ {
			Primt(Coordinates{
				X: ballCorner.X + j,
				Y: ballCorner.Y + i,
			}, " ")
		}
	}
}

func (ballstate *Ball) RenderBall() {

	var ballCorner = ballstate.GetCornerPointer("TL")

	for i := float64(0); i < ballstate.Height; i++ {
		fmt.Print("\033[" + fmt.Sprint(ballCorner.Y+i) + ";" + fmt.Sprint(ballCorner.X) + "H")
		for j := float64(0); j < ballstate.Width; j++ {
			if i == float64(0) || i == ballstate.Height-1 {
				if j == float64(0) || j == ballstate.Width-1 {
					Primt(Coordinates{
						X: ballCorner.X + j,
						Y: ballCorner.Y + i,
					}, " ")
					continue
				}
			}
			Primt(Coordinates{
				X: ballCorner.X + j,
				Y: ballCorner.Y + i,
			}, "#")
		}
	}
}

func (ballstate *Ball) MoveBall() {
	ballstate.ClearBall()

	uvx, uvy := ballstate.Velocity.GetUnitVector()
	//fmt.Println(ballstate.Position.X, ballstate.Position.Y)

	ballstate.Position.X += uvx
	ballstate.Position.Y += uvy

	//fmt.Println(uvx, uvy)
	//fmt.Println(ballstate.Position.X, ballstate.Position.Y)

	ballstate.RenderBall()
}
