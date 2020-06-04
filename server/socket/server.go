package socket

import (
	"context"
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

	loginID := request.LoginID

	if s.config.VerifyUser {
		panic("not implemented")
	}

	player, err := s.games.CreatePlayer(loginID, loginID)
	if err != nil {
		conn.sendError("invalid login", false)
		return nil, err
	}

	conn.playerID = player.UUID
	conn.sendMessage(LoginSuccess{
		Username: player.Username,
		PlayerID: player.UUID,
	})

	return &conn, nil
}

func (s *Server) runMessageLoop(conn *clientConn) {
	incomingMsgs := s.listenForMessages(conn)

	openGames, newGameCh, err := s.games.OpenGamesForPlayer(conn.playerID)
	if err != nil {
		panic(err)
	}
	defer s.games.CloseGames(openGames)

	for _, g := range openGames {
		s.handleGameUpdate(conn, g.Game)
	}
	var newestGame store.NewGameNotification
	defer newestGame.Game.Close(newestGame.UpdateCh)

	ctx, cancelCtx := context.WithCancel(context.TODO())
	conn.cancelCtx = cancelCtx

	openGamesCh := s.games.ListenAny(openGames, ctx)
loop:
	for {
		select {
		case msg, ok := <-incomingMsgs:
			if !ok {
				// socket closed, end message loop
				break loop
			}

			s.handleMessage(conn, msg, openGames)
			break

		// either a game update or we need to update our new game
		case gameIdx, ok := <-openGamesCh:
			if !ok {
				if newestGame.Game != nil {
					openGames = append(openGames, newestGame)
				}
				ctx, cancelCtx = context.WithCancel(context.TODO())
				openGamesCh = s.games.ListenAny(openGames, ctx)
				newestGame = store.NewGameNotification{} // zero out old struct
				break
			}
			s.handleGameUpdate(conn, openGames[gameIdx].Game)
			break
		case newestGame = <-newGameCh:
			s.handleGameUpdate(conn, newestGame.Game)
			cancelCtx()
			break
		}
	}
}

func (s *Server) handleGameUpdate(conn *clientConn, g *store.Game) {
	state, err := g.GetGameState(conn.playerID)
	if err != nil {
		panic(err)
	}

	conn.sendMessage(state)
}

func (s *Server) handleMessage(conn *clientConn, msg interface{}, games []store.NewGameNotification) {
	switch v := msg.(type) {
	case *NewGame:
		s.handleNewGame(conn, v)
		break
	case *PlayMove:
		idx := 0
		for idx = range games {
			if games[idx].Game.UUID() == v.GameID {
				err := games[idx].Game.PlayMove(v.Move)
				if err != nil {
					panic(err)
				}
			}
		}
		break
	case *UserLookup:
		username := v.Username
		playerID := v.PlayerID
		if username != "" {
			s.handleLookupByUsername(conn, username)
		} else if playerID != "" {
			panic("not impl")
		} else {
			conn.sendError("Must specify either username or playerID", true)
		}

	default:
		fmt.Printf("Unknown message type %+v\n", msg)
	}
}

func (s *Server) handleLookupByUsername(conn *clientConn, username string) {
	fullplayer, err := s.games.TryLookupPlayerUsername(username)
	if err != nil {
		conn.sendError("Could not lookup user", true)
		return
	}

	err = conn.sendMessage(UserLookup{
		Username: fullplayer.Username,
		PlayerID: fullplayer.UUID,
	})
	if err != nil {
		panic(err)
	}
}

func (s *Server) handleNewGame(conn *clientConn, payload *NewGame) {
	err := s.games.NewGame(conn.playerID, payload.OpponentID)
	if err != nil {
		panic(err)
	}
}

func (s *Server) listenForMessages(conn *clientConn) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer recover() // TODO: Log errored socket
		defer close(ch)
		for {
			msg, err := conn.nextMessageSync()
			if err != nil {
				break // TODO: Log errored socket
			}

			ch <- msg
		}
	}()
	return ch
}
