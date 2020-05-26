package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	*config
	upgrader websocket.Upgrader
}

func NewSocketServer(c *config) *SocketServer {
	checkOriginFunc := func(*http.Request) bool {
		return true
	}

	if c.CheckOrigin {
		// use default origin checker
		checkOriginFunc = nil
	}
	return &SocketServer{
		c,
		websocket.Upgrader{
			ReadBufferSize:  512,
			WriteBufferSize: 1024,
			CheckOrigin:     checkOriginFunc,
		},
	}
}

func (s *SocketServer) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	playerID, err := s.login(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = playerID
}

func (s *SocketServer) login(conn *websocket.Conn) (string, error) {
	_, reader, err := conn.NextReader()
	if err != nil {
		return "", err
	}

	req, err := parseMessage(reader)
	request, ok := req.(*LoginRequest)
	if err != nil {
		return "", err
	} else if !ok {
		return "", fmt.Errorf("wrong type recieved: %#v", req)
	}

	if s.config.VerifyUser {
		panic("not implemented")
	}

	sendMessage(conn, LoginSuccess{})

	return request.PlayerID, nil
}

func sendMessage(conn *websocket.Conn, payload interface{}) error {
	return conn.WriteJSON(OutgoingSocketMessage{
		Type:    reflect.TypeOf(payload).Name(),
		Payload: payload,
	})
}

func parseMessage(r io.Reader) (interface{}, error) {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var message IncomingSocketMessage
	err := decoder.Decode(&message)
	if err != nil {
		return nil, err
	}

	decodedMessage := valueFromType(message.Type)

	err = json.Unmarshal(message.Payload, decodedMessage)
	if err != nil {
		return nil, err
	}

	return decodedMessage, nil
}

func valueFromType(typ string) interface{} {
	switch typ {
	case "LoginRequest":
		return &LoginRequest{}
	}

	return nil
}

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
