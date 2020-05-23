package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	*config
	upgrader websocket.Upgrader
}

func NewSocketServer(c *config) *SocketServer {
	return &SocketServer{
		c,
		websocket.Upgrader{
			ReadBufferSize:  512,
			WriteBufferSize: 1024,
		},
	}
}

func (s *SocketServer) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	playerID, err := s.login(conn)
	if err != nil {
		// TODO: write some error message
		conn.Close()
		return
	}

	_ = playerID
	// TODO: handle other messages

	conn.WriteJSON("test")

	conn.Close()
}

func (s *SocketServer) login(conn *websocket.Conn) (string, error) {
	_, reader, err := conn.NextReader()
	if err != nil {
		return "", err
	}

	request := LoginRequest{}
	err = json.NewDecoder(reader).Decode(&request)
	if err != nil {
		return "", err
	}

	if s.config.VerifyUser {
		panic("not implemented")
	} else {
		return request.PlayerID, nil
	}
}

type LoginRequest struct {
	PlayerID      string
	Authorization string
}
