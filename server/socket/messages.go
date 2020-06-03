package socket

import "encoding/json"

type IncomingSocketMessage struct {
	Type    string          `json:"messageType"`
	Payload json.RawMessage `json:"payload"`
}

type OutgoingSocketMessage struct {
	Type    string      `json:"messageType"`
	Payload interface{} `json:"payload"`
}

type LoginRequest struct {
	PlayerID      string `json:"playerID"`
	Authorization string `json:"authorization"`
}

type LoginSuccess struct {
}

type ErrorMessage struct {
	Message string `json:"message"`

	// if Recoverable is false, then the websocket is closed after this
	// message is sent
	Recoverable bool `json:"recoverable"`
}
