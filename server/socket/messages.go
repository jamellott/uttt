package socket

import (
	"encoding/json"

	"github.com/heartles/uttt/server/game"
)

type IncomingSocketMessage struct {
	Type    string          `json:"messageType"`
	Payload json.RawMessage `json:"payload"`
}

type OutgoingSocketMessage struct {
	Type    string      `json:"messageType"`
	Payload interface{} `json:"payload"`
}

type LoginRequest struct {
	LoginID string `json:"loginID"`
}

type NewGame struct {
	Opponent string `json:"opponent"`
}

type PlayMove struct {
	GameID string    `json:"gameID"`
	Move   game.Move `json:"move"`
}

type GameUpdate struct {
	GameID     string      `json:"gameID"`
	PlayerX    string      `json:"playerX"`
	PlayerO    string      `json:"playerO"`
	ValidMoves []game.Move `json:"validMoves"`
	Victor     string      `json:"victor"`
}

type LoginSuccess struct {
	Username string       `json:"username"`
	PlayerID string       `json:"playerID"`
	Games    []GameUpdate `json:"games"`
}

type ErrorMessage struct {
	Message string `json:"message"`

	// if Recoverable is false, then the websocket is closed after this
	// message is sent
	Recoverable bool `json:"recoverable"`
}
