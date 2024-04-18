package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsMessage struct {
	MsgType   string  `json:"Type"`
	BallState Ball    `json:"BallState,omitempty"`
	CtrlMsg   CtrlMsg `json:"CtrlMsg,omitempty"`
}

type Ball struct {
	Height, Width float64
	Position      Coordinates
	Velocity      Velocity
}
type Coordinates struct {
	X, Y float64
}
type Velocity struct {
	X, Y int
}
type CtrlMsg struct {
	ConnectionState string      `json:"ConnectionState"`
	Payload         interface{} `json:"Payload"`
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("Server started on :8081/ws")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Println("Error starting server:", err)
		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		_ = SendAndReturnError(conn, "Error upgrading connection"+err.Error())
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			_ = SendAndReturnError(conn, "Error closing connection"+err.Error())
		}
	}()
	err = HandleINITmsg(conn)
	if err != nil {
		if err != nil {
			log.Println("Error writing close message:", err)
			return
		}
		return
	}
	for {
		var msg wsMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			_ = SendAndReturnError(conn, "Error reading JSON message")
			return
		}
		if msg.MsgType == "CTRL" {
			HandleCTRLmsg(conn, msg)
		} else if msg.MsgType == "BALL" {
			HandleBALLmsg(conn, msg)
		} else {
			_ = SendAndReturnError(conn, "Invalid message type:"+msg.MsgType)
			return
		}
	}
}

func HandleINITmsg(conn *websocket.Conn) error {
	initMessage := wsMessage{}
	err := conn.ReadJSON(&initMessage)
	if err != nil {
		return SendAndReturnError(conn, "Error reading JSON message")
	}
	if initMessage.MsgType != "CTRL" {
		return SendAndReturnError(conn, "Invalid message type:"+initMessage.MsgType)
	}
	if initMessage.CtrlMsg.ConnectionState != "INIT" {
		return SendAndReturnError(conn, "Unexpected connection state: "+initMessage.CtrlMsg.ConnectionState)
	}
	if v, ok := initMessage.CtrlMsg.Payload.(string); !ok || v == "" {
		return SendAndReturnError(conn, "Invalid Username")
	}
	username := initMessage.CtrlMsg.Payload.(string)
	CurrOnlineUsers := GetOnlineUsers()
	listOfOnlineUsers := make([]string, 0)
	for user, _ := range CurrOnlineUsers {
		listOfOnlineUsers = append(listOfOnlineUsers, user)
	}
	ok := RegisterUser(username, conn)
	if !ok {
		return SendAndReturnError(conn, "Error registering user")
	}
	resMessage := wsMessage{
		MsgType: "CTRL",
		CtrlMsg: CtrlMsg{
			ConnectionState: "ACK",
			Payload:         listOfOnlineUsers,
		},
	}
	err = conn.WriteJSON(resMessage)
	if err != nil {
		return SendAndReturnError(conn, "Error Sending ACK message"+err.Error())
	}

	return nil
}

func HandleCTRLmsg(conn *websocket.Conn, msg wsMessage) {
	switch msg.CtrlMsg.ConnectionState {
	case "CONNECT":
		{
			if v, ok := msg.CtrlMsg.Payload.(string); !ok || v == "" {
				err := conn.WriteJSON(wsMessage{
					MsgType: "CTRL",
					CtrlMsg: CtrlMsg{
						ConnectionState: "GAME_DEN",
						Payload:         "Invalid Username",
					},
				})
				if err != nil {
					log.Println("Error writing JSON message:", err)
					return
				}
				return
			}

			username := msg.CtrlMsg.Payload.(string)

			// send startgame request
			otherUserConn := GetOnlineUsers()[username]
			if otherUserConn == nil {
				err := conn.WriteJSON(wsMessage{
					MsgType: "CTRL",
					CtrlMsg: CtrlMsg{
						ConnectionState: "GAME_DEN",
						Payload:         "User not online",
					},
				})
				if err != nil {
					log.Println("Error writing JSON message:", err)
					return
				}
				return
			}

			err := otherUserConn.WriteJSON(wsMessage{
				MsgType: "CTRL",
				CtrlMsg: CtrlMsg{
					ConnectionState: "CONN_REQ",
					Payload:         OnlineUsers.List[conn].Username,
				},
			})
			if err != nil {
				_ = SendAndReturnError(conn, "Error sending connection request")
				return
			}
		}
	case "CONN_REJ":
		{
			if v, ok := msg.CtrlMsg.Payload.(string); ok && v != "" {
				otherUserConn, ok := GetOnlineUsers()[v]
				if ok {
					err := otherUserConn.WriteJSON(wsMessage{
						MsgType: "CTRL",
						CtrlMsg: CtrlMsg{
							ConnectionState: "GAME_DEN",
							Payload:         "User rejected connection",
						},
					})
					if err != nil {
						log.Println("Error writing JSON message:", err)
						return
					}
				}
			}
		}
	case "CONN_ACC":
		{
			if v, ok := msg.CtrlMsg.Payload.(string); !ok || v == "" {
				err := conn.WriteJSON(wsMessage{
					MsgType: "CTRL",
					CtrlMsg: CtrlMsg{
						ConnectionState: "GAME_DEN",
						Payload:         "Invalid Username",
					},
				})
				if err != nil {
					log.Println("Error writing JSON message:", err)
					return
				}
				return
			}
			otherUserConn := GetOnlineUsers()[msg.CtrlMsg.Payload.(string)]
			if otherUserConn == nil {
				err := conn.WriteJSON(wsMessage{
					MsgType: "CTRL",
					CtrlMsg: CtrlMsg{
						ConnectionState: "GAME_DEN",
						Payload:         "User not online",
					},
				})
				if err != nil {
					log.Println("Error writing JSON message:", err)
					return
				}
				return
			}
			err := StartGame(conn, otherUserConn)
			if err != nil {
				log.Println("Error starting game:", err)
				return
			}
		}
	case "GAME_OVER":
		{

			otherUserConn := OnlineUsers.List[conn].Apponent

			if otherUserConn != nil {
				err := otherUserConn.WriteJSON(wsMessage{
					MsgType: "CTRL",
					CtrlMsg: CtrlMsg{
						ConnectionState: "GAME_OVER",
						Payload:         "You Won!",
					},
				})
				if err != nil {
					log.Println("Error writing JSON message:", err)
					return
				}
			}

			OnlineUsers.Lock()
			OnlineUsers.List[conn].Playing = false
			OnlineUsers.List[conn].Apponent = nil
			OnlineUsers.List[otherUserConn].Apponent = nil
			OnlineUsers.List[otherUserConn].Playing = false
			OnlineUsers.Unlock()
		}
	}
}

func HandleBALLmsg(conn *websocket.Conn, msg wsMessage) {
	otherUserConn := OnlineUsers.List[conn].Apponent
	invertedBall := Ball{
		Position: Coordinates{
			X: TERMINAL_WIDTH - msg.BallState.Position.X,
			Y: msg.BallState.Position.Y,
		},
		Velocity: Velocity{
			X: -msg.BallState.Velocity.X,
			Y: -msg.BallState.Velocity.Y,
		},
	}
	err := otherUserConn.WriteJSON(wsMessage{
		MsgType:   "BALL",
		BallState: invertedBall,
	})
	if err != nil {
		log.Println("Error writing JSON message:", err)
		return
	}
}

func SendAndReturnError(conn *websocket.Conn, errMsg string) error {
	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, errMsg))
	if err != nil {
		errors.Join(errors.New("Error writing close message:"), err)
		return err
	}
	return err
}
