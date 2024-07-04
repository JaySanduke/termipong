package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type wsMessage struct {
	MsgType   string  `json:"Type"`
	BallState Ball    `json:"BallState,omitempty"`
	CtrlMsg   CtrlMsg `json:"CtrlMsg,omitempty"`
}

type CtrlMsg struct {
	ConnectionState string      `json:"ConnectionState"`
	Payload         interface{} `json:"Payload"`
}

var addr = flag.String("server", "10.0.0.140:8081", "http service address")

func ConnectionInit() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	//TODO handle err

	return c
}

func SendMessage(message wsMessage, conn *websocket.Conn) {
	err := conn.WriteJSON(message)
	if err != nil {
		log.Println("write:", err)
	}
}
