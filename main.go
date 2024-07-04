package main

import (
	"fmt"
	"log"
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

var fpsController = time.Tick(50 * time.Millisecond)

func InitGame() {
	var newBall = Ball{
		Height:   3,
		Width:    5,
		Position: Coordinates{X: 10, Y: 10},
		Velocity: Velocity{X: 5, Y: 3},
	}

	var newboard = Board{Height: 2, Width: 10, Position: Coordinates{X: 5, Y: 5}}

	newboard.RenderBoard(&TerminalSize)
	//newBall.RenderBall()
	func() {
		isGameRunning := true

		for isGameRunning {
			select {
			case <-fpsController:
				newBall.MoveBall()
				isGameRunning = CollisionDetection(&newBall, &newboard)
			}
		}
	}()
}

func WelcomeText(steadystate chan bool) {

	var textarray = []string{
		"$$$$$$$$\\       $$$$$$$$\\       $$$$$$$\\        $$\\      $$\\       $$$$$$\\       $$$$$$$\\         $$$$$$\\        $$\\   $$\\        $$$$$$\\  \n\\__$$  __|      $$  _____|      $$  __$$\\       $$$\\    $$$ |      \\_$$  _|      $$  __$$\\       $$  __$$\\       $$$\\  $$ |      $$  __$$\\ \n   $$ |         $$ |            $$ |  $$ |      $$$$\\  $$$$ |        $$ |        $$ |  $$ |      $$ /  $$ |      $$$$\\ $$ |      $$ /  \\__|\n   $$ |         $$$$$\\          $$$$$$$  |      $$\\$$\\$$ $$ |        $$ |        $$$$$$$  |      $$ |  $$ |      $$ $$\\$$ |      $$ |$$$$\\ \n   $$ |         $$  __|         $$  __$$<       $$ \\$$$  $$ |        $$ |        $$  ____/       $$ |  $$ |      $$ \\$$$$ |      $$ |\\_$$ |\n   $$ |         $$ |            $$ |  $$ |      $$ |\\$  /$$ |        $$ |        $$ |            $$ |  $$ |      $$ |\\$$$ |      $$ |  $$ |\n   $$ |         $$$$$$$$\\       $$ |  $$ |      $$ | \\_/ $$ |      $$$$$$\\       $$ |             $$$$$$  |      $$ | \\$$ |      \\$$$$$$  |\n   \\__|         \\________|      \\__|  \\__|      \\__|     \\__|      \\______|      \\__|             \\______/       \\__|  \\__|       \\______/ \n                                                                                                                                           \n                                                                                                                                           \n                                                                                                                                           ",
		" /$$$$$$$$       /$$$$$$$$       /$$$$$$$        /$$      /$$       /$$$$$$       /$$$$$$$         /$$$$$$        /$$   /$$        /$$$$$$ \n|__  $$__/      | $$_____/      | $$__  $$      | $$$    /$$$      |_  $$_/      | $$__  $$       /$$__  $$      | $$$ | $$       /$$__  $$\n   | $$         | $$            | $$  \\ $$      | $$$$  /$$$$        | $$        | $$  \\ $$      | $$  \\ $$      | $$$$| $$      | $$  \\__/\n   | $$         | $$$$$         | $$$$$$$/      | $$ $$/$$ $$        | $$        | $$$$$$$/      | $$  | $$      | $$ $$ $$      | $$ /$$$$\n   | $$         | $$__/         | $$__  $$      | $$  $$$| $$        | $$        | $$____/       | $$  | $$      | $$  $$$$      | $$|_  $$\n   | $$         | $$            | $$  \\ $$      | $$\\  $ | $$        | $$        | $$            | $$  | $$      | $$\\  $$$      | $$  \\ $$\n   | $$         | $$$$$$$$      | $$  | $$      | $$ \\/  | $$       /$$$$$$      | $$            |  $$$$$$/      | $$ \\  $$      |  $$$$$$/\n   |__/         |________/      |__/  |__/      |__/     |__/      |______/      |__/             \\______/       |__/  \\__/       \\______/ \n                                                                                                                                           \n                                                                                                                                           \n                                                                                                                                           ",
		" ________        ________        _______         __       __        ______        _______          ______         __    __         ______  \n|        \\      |        \\      |       \\       |  \\     /  \\      |      \\      |       \\        /      \\       |  \\  |  \\       /      \\ \n \\$$$$$$$$      | $$$$$$$$      | $$$$$$$\\      | $$\\   /  $$       \\$$$$$$      | $$$$$$$\\      |  $$$$$$\\      | $$\\ | $$      |  $$$$$$\\\n   | $$         | $$__          | $$__| $$      | $$$\\ /  $$$        | $$        | $$__/ $$      | $$  | $$      | $$$\\| $$      | $$ __\\$$\n   | $$         | $$  \\         | $$    $$      | $$$$\\  $$$$        | $$        | $$    $$      | $$  | $$      | $$$$\\ $$      | $$|    \\\n   | $$         | $$$$$         | $$$$$$$\\      | $$\\$$ $$ $$        | $$        | $$$$$$$       | $$  | $$      | $$\\$$ $$      | $$ \\$$$$\n   | $$         | $$_____       | $$  | $$      | $$ \\$$$| $$       _| $$_       | $$            | $$__/ $$      | $$ \\$$$$      | $$__| $$\n   | $$         | $$     \\      | $$  | $$      | $$  \\$ | $$      |   $$ \\      | $$             \\$$    $$      | $$  \\$$$       \\$$    $$\n    \\$$          \\$$$$$$$$       \\$$   \\$$       \\$$      \\$$       \\$$$$$$       \\$$              \\$$$$$$        \\$$   \\$$        \\$$$$$$ \n                                                                                                                                           \n                                                                                                                                           ",
		" ________        ________        _______         __       __        ______        _______          ______         __    __         ______  \n/        |      /        |      /       \\       /  \\     /  |      /      |      /       \\        /      \\       /  \\  /  |       /      \\ \n$$$$$$$$/       $$$$$$$$/       $$$$$$$  |      $$  \\   /$$ |      $$$$$$/       $$$$$$$  |      /$$$$$$  |      $$  \\ $$ |      /$$$$$$  |\n   $$ |         $$ |__          $$ |__$$ |      $$$  \\ /$$$ |        $$ |        $$ |__$$ |      $$ |  $$ |      $$$  \\$$ |      $$ | _$$/ \n   $$ |         $$    |         $$    $$<       $$$$  /$$$$ |        $$ |        $$    $$/       $$ |  $$ |      $$$$  $$ |      $$ |/    |\n   $$ |         $$$$$/          $$$$$$$  |      $$ $$ $$/$$ |        $$ |        $$$$$$$/        $$ |  $$ |      $$ $$ $$ |      $$ |$$$$ |\n   $$ |         $$ |_____       $$ |  $$ |      $$ |$$$/ $$ |       _$$ |_       $$ |            $$ \\__$$ |      $$ |$$$$ |      $$ \\__$$ |\n   $$ |         $$       |      $$ |  $$ |      $$ | $/  $$ |      / $$   |      $$ |            $$    $$/       $$ | $$$ |      $$    $$/ \n   $$/          $$$$$$$$/       $$/   $$/       $$/      $$/       $$$$$$/       $$/              $$$$$$/        $$/   $$/        $$$$$$/  \n                                                                                                                                           \n                                                                                                                                           ",
	}

	var num = 1

	var count = 0
	for {
		num = count % 4
		fmt.Print(textarray[num] + "\n")
		count++

		select {
		case <-steadystate:
			print("\033[2J\033[0;0H")
			return
		default:
			fmt.Println("\nWaiting for server to connect...")
			time.Sleep(260 * time.Millisecond)
			print("\033[2J\033[0;0H")
		}

	}
}

//func GetUsername() string {
//
//	return user
//}

func ShowHomePage() {
	welcometextstate := make(chan bool)

	go WelcomeText(welcometextstate)

	//onlineusers := make(chan []string)

	//TODO websocket conn init

	// connection estb
	websocketconn := ConnectionInit()

	welcometextstate <- true

	//fmt.Println("waiting")
	//time.Sleep(10 * time.Second)

	// username enter
	var username string
	//var user string
	fmt.Print("Enter Username: ")
	n, err := fmt.Scanln(&username)
	if err != nil {
		log.Fatal(err)
	}

	if n <= 0 {
		fmt.Println("Invalid input")
	}

	if username == "" {
		fmt.Println("Please enter a username")
	}

	SendMessage(wsMessage{
		MsgType: "CTRL",
		CtrlMsg: CtrlMsg{
			ConnectionState: "INIT",
			Payload:         username,
		},
	}, websocketconn)
	// list of online users

	onlineUsers := wsMessage{}

	err = websocketconn.ReadJSON(&onlineUsers)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(onlineUsers)

	// select opponent

	//time.Sleep(2 * time.Second)
	//welcometextstate <- true

	TheGameMode = "GAME"
	fmt.Scanln()
	return

}

func ShowGamePage() {
	//print("game page")
	InitGame()
}

func GameMode(mode string) {
	switch mode {
	case "HOME":
		ShowHomePage()
	case "GAME":
		ShowGamePage()
	}
}

var (
	TheGameMode string = "HOME" //default game state i.e. :  HOME, GAME
)

func main() {
	print("\033[2J\033[0;0H")

	for {
		GameMode(TheGameMode)
	}

	fmt.Scanln()
}
