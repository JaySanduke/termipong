package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"sync"
)

type UsersList struct {
	List map[*websocket.Conn]*User
	sync.Mutex
}
type ConnList struct {
	List map[string]*websocket.Conn
	sync.Mutex
}

const TERMINAL_HEIGHT = 40
const TERMINAL_WIDTH = 160

var (
	OnlineUsers = UsersList{
		List: make(map[*websocket.Conn]*User),
	}
	OnlineConns = ConnList{
		List: make(map[string]*websocket.Conn),
	}
)

type User struct {
	Username string
	Conn     *websocket.Conn
	Playing  bool
	Apponent *websocket.Conn
}

func GetOnlineUsers() map[string]*websocket.Conn {
	return OnlineConns.List
}
func RegisterUser(username string, conn *websocket.Conn) bool {

	if _, ok := OnlineConns.List[username]; ok {
		_ = SendAndReturnError(conn, "Username already taken")
		return false
	}
	OnlineUsers.Lock()
	defer OnlineUsers.Unlock()
	OnlineConns.Lock()
	defer OnlineConns.Unlock()
	OnlineUsers.List[conn] = &User{
		Username: username,
		Conn:     conn,
		Playing:  false,
	}
	OnlineConns.List[username] = conn
	return true
}

func UnregisterUser(conn *websocket.Conn) {
	OnlineUsers.Lock()
	OnlineConns.Lock()
	defer OnlineUsers.Unlock()
	defer OnlineConns.Unlock()
	delete(OnlineConns.List, OnlineUsers.List[conn].Username)
	delete(OnlineUsers.List, conn)
}

func StartGame(P1 *websocket.Conn, P2 *websocket.Conn) error {
	OnlineUsers.Lock()
	OnlineUsers.List[P1].Playing = true
	OnlineUsers.List[P2].Playing = true
	OnlineUsers.List[P1].Apponent = P2
	OnlineUsers.List[P2].Apponent = P1
	OnlineUsers.Unlock()

	noBallmsg := wsMessage{
		MsgType: "CTRL",
		CtrlMsg: CtrlMsg{
			ConnectionState: "GAME_START",
			Payload:         fmt.Sprint(TERMINAL_WIDTH, "x", TERMINAL_HEIGHT),
		},
		BallState: Ball{
			Height:   3,
			Width:    5,
			Position: Coordinates{-10, -10},
			Velocity: Velocity{0, 0},
		},
	}
	ballMsg := wsMessage{
		MsgType: "CTRL",
		CtrlMsg: CtrlMsg{
			ConnectionState: "GAME_START",
			Payload:         fmt.Sprint(TERMINAL_WIDTH, "x", TERMINAL_HEIGHT),
		},
		BallState: GetRandomBallState(),
	}
	if ballOwner := rand.Int() % 2; ballOwner == 0 {
		err := P1.WriteJSON(ballMsg)
		if err != nil {
			return err
		}
		err = P2.WriteJSON(noBallmsg)
		if err != nil {
			return err
		}
	} else {
		err := P2.WriteJSON(ballMsg)
		if err != nil {
			return err
		}
		err = P1.WriteJSON(noBallmsg)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetRandomBallState() Ball {
	rInt := rand.Int()
	ballXPos := 50 + float64(rInt%60)

	rInt = rand.Int()
	ballYPos := 2 + float64(rInt%5)

	rInt = rand.Int()
	ballXVel := 10 - float64(rInt%20)

	if ballXVel == 0 {
		ballXVel = 1
	}

	rInt = rand.Int()
	ballYVel := 1 + float64(rInt%10)

	return Ball{
		Height:   3,
		Width:    5,
		Position: Coordinates{X: ballXPos, Y: ballYPos},
		Velocity: Velocity{X: int(ballXVel), Y: int(ballYVel)},
	}
}
