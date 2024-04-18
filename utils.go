package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
)

// utils
func GetTerminalSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 80, 24 // default size if unable to get terminal size
	}
	var width, height int
	fmt.Sscanf(string(out), "%d %d", &height, &width)
	return width, height
}

func Primt(Position Coordinates, text string) {

	if Position.Y < 0 || Position.X < 0 {
		return
	}
	if Position.X >= float64(TerminalSize.width) || Position.Y >= float64(TerminalSize.height) {
		return
	}

	fmt.Print("\033[" + fmt.Sprint(int(math.Round(Position.Y))) + ";" + fmt.Sprint(int(math.Round(Position.X))) + "H")
	fmt.Print(text)
}
