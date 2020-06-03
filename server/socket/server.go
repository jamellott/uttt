package socket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/heartles/uttt/server/config"
	"github.com/heartles/uttt/server/store"
)

type Server struct {
	config   *config.Config
	upgrader websocket.Upgrader
	games    *store.GameService
}

func NewServer(c *config.Config, gameSvc *store.GameService) *Server {
	checkOriginFunc := func(*http.Request) bool {
		return true
	}

	if c.CheckOrigin {
		// use default origin checker
		checkOriginFunc = nil
	}
	return &Server{
		c,
		websocket.Upgrader{
			ReadBufferSize:  512,
			WriteBufferSize: 1024,
			CheckOrigin:     checkOriginFunc,
		},
		gameSvc,
	}
}

func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	socket, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn, err := s.login(socket)
	if err != nil {
		fmt.Println(err)
		socket.Close()
		return
	}

	s.runMessageLoop(conn)
}

func (s *Server) login(socket *websocket.Conn) (*clientConn, error) {
	conn := clientConn{
		socket: socket,
	}

	req, err := conn.nextMessageSync()
	request, ok := req.(*LoginRequest)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("wrong type recieved: %#v", req)
	}

	if s.config.VerifyUser {
		panic("not implemented")
	}

	conn.playerID = request.PlayerID
	player, err := s.games.CreatePlayer(conn.playerID, conn.playerID)
	if err != nil {
		conn.sendError("invalid login", false)
		return nil, err
	}

	conn.sendMessage(LoginSuccess{
		Username: player.Username,
		PlayerID: player.UUID,
	})

	return &conn, nil
}

func (s *Server) runMessageLoop(conn *clientConn) {
	incomingMsgs := s.listenForMessages(conn)

loop:
	for {
		select {
		case msg, ok := <-incomingMsgs:
			if !ok {
				// socket closed, end message loop
				break loop
			}

			s.handleMessage(conn, msg)
			break
			// TODO: add case for "other player played move/sent chat"
		}
	}
}

func (s *Server) handleMessage(conn *clientConn, msg interface{}) {
	switch msg.(type) {
	}
}

func (s *Server) listenForMessages(conn *clientConn) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for {
			msg, err := conn.nextMessageSync()
			if err == nil {
				break
			}

			ch <- msg
		}
		close(ch)
	}()
	return ch
}
