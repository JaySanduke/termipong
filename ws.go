package main

type wsMessage struct {
	MsgType   string  `json:"Type"`
	BallState Ball    `json:"BallState,omitempty"`
	CtrlMsg   CtrlMsg `json:"CtrlMsg,omitempty"`
}

type CtrlMsg struct {
	ConnectionState string      `json:"ConnectionState"`
	Payload         interface{} `json:"Payload"`
}
