package store

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/heartles/uttt/server/game"
)

// TODO: There's a ton of race conditions in here

type Game struct {
	underlying     *game.Game
	mutex          sync.RWMutex
	service        *GameService
	uuid           string
	listenChannels []chan struct{}
}

func (g *Game) UUID() string {
	return g.uuid
}

func (g *Game) Close(ch <-chan struct{}) error {
	if g == nil {
		return nil
	}
	g.mutex.Lock()

	idx := 0
	found := false
	for idx = range g.listenChannels {
		if g.listenChannels[idx] == ch {
			found = true
			break
		}
	}
	if !found {
		g.mutex.Unlock()
		return errors.New("Invalid channel")
	}

	close(g.listenChannels[idx])
	g.listenChannels = append(g.listenChannels[:idx], g.listenChannels[idx:]...)

	g.mutex.Unlock()
	return g.service.closeGame(g)
}

type SquareState struct {
	Owner      *string         `json:"owner"`
	Playable   bool            `json:"playable"`
	Coordinate game.Coordinate `json:"coordinate"`
}

type GridState struct {
	Owner   *string           `json:"owner"`
	Squares [3][3]SquareState `json:"squares"`
}

type GameState struct {
	GameID  string `json:"gameID"`
	PlayerX string `json:"playerX"`
	PlayerO string `json:"playerO"`

	PlayerXName string `json:"playerXName"`
	PlayerOName string `json:"playerOName"`

	Victor *string         `json:"victor"`
	Grids  [3][3]GridState `json:"grids"`
}

func (g *Game) GetGameState(playerID string) (*GameState, error) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	validMoves := g.underlying.GetValidMoves(playerID)

	isPlayable := func(coord game.Coordinate) bool {
		for i := range validMoves {
			if validMoves[i].Coordinate == coord {
				return true
			}
		}
		return false
	}

	playerX, playerO, _, _ := g.underlying.SaveGame()
	playerXFull, _ := g.service.TryLookupPlayerUUID(playerX)
	playerOFull, _ := g.service.TryLookupPlayerUUID(playerO)
	gameState := &GameState{
		GameID:      g.uuid,
		PlayerX:     playerX,
		PlayerO:     playerO,
		PlayerXName: playerXFull.Username,
		PlayerOName: playerOFull.Username,
	}

	victor := g.underlying.GameWinner()
	if victor != "" {
		gameState.Victor = &victor
	}
	// GameCoordinate = {{z, w} {x, y}}
	for z := 0; z <= 2; z++ {
		for w := 0; w <= 2; w++ {
			grid := &gameState.Grids[w][z]

			owner, _ := g.underlying.BlockWinner(game.SubCoordinate{z, w})

			if owner != "" {
				grid.Owner = &owner
			}

			for x := 0; x <= 2; x++ {
				for y := 0; y <= 2; y++ {
					coord := game.NewCoordinate(z+1, w+1, x+1, y+1)
					square := &grid.Squares[y][x]
					owner, _ := g.underlying.SquareOwner(coord)
					if owner != "" {
						square.Owner = &owner
					}
					square.Playable = isPlayable(coord)
					square.Coordinate = coord
				}
			}
		}
	}

	return gameState, nil
}

func (g *Game) PlayMove(m game.Move) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return g.underlying.PlayMove(m)
}

func (g *Game) listenForUpdates() <-chan struct{} {
	unbufferedCh := make(chan struct{})
	ch := make(chan struct{})
	// avoid potential (but unlikely) deadlock
	go func() {
		var updateOccurred bool

		for {
			defer close(ch)
			if updateOccurred {
				select {
				case ch <- struct{}{}:
					break

				case _, ok := <-unbufferedCh:
					if !ok {
						return
					}

					updateOccurred = true
					break
				}
			} else {
				_, ok := <-unbufferedCh
				if !ok {
					return
				}

				updateOccurred = true
			}
		}
	}()
	g.mutex.Lock()
	defer g.mutex.Unlock()

	return ch
}

type GameService struct {
	games   map[string]*loadedGame
	players map[string]chan NewGameNotification
	mutex   sync.Mutex
	*Store
}

type NewGameNotification struct {
	Game     *Game
	UpdateCh <-chan struct{}
}

func NewGameService(dbFilename string) (*GameService, error) {
	st, err := NewStore(dbFilename)
	if err != nil {
		return nil, err
	}

	return &GameService{
		map[string]*loadedGame{},
		map[string]chan NewGameNotification{},
		sync.Mutex{},
		st,
	}, nil
}

type loadedGame struct {
	openConns int
	game      *Game
}

func (s *GameService) closeGame(g *Game) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	loaded := s.games[g.uuid]

	loaded.openConns--
	if loaded.openConns == 0 {
		// unload game
		delete(s.games, g.uuid)
		return s.Store.saveGame(g.uuid, g.underlying)
	}

	return nil
}

func (s *GameService) ListenAny(notifs []NewGameNotification, ctx context.Context) <-chan int {
	if len(notifs) == 0 {
		return nil
	}

	ch := make(chan int)

	go func() {
		defer close(ch)
		allChs := make([]reflect.SelectCase, len(notifs)+1)
		for {
			for i := range notifs {
				allChs[i] = reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(notifs[i].UpdateCh),
				}
			}
			allChs[len(notifs)] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ctx.Done()),
			}
			chosen, _, _ := reflect.Select(allChs)

			// ctx canceled
			if chosen == len(notifs) {
				return
			}
			ch <- chosen
		}
	}()

	return ch
}

func (s *GameService) CloseGames(notifs []NewGameNotification) {
	for _, n := range notifs {
		n.Game.Close(n.UpdateCh)
	}
}

func (s *GameService) CloseNewGameCh(playerID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	ch := s.players[playerID]
	close(ch)

	delete(s.players, playerID)
}

func (s *GameService) NewGame(playerX string, playerO string) error {
	g, err := game.NewGame(playerX, playerO)
	if err != nil {
		return err
	}

	uuid, err := s.Store.saveNewGame(g)
	if err != nil {
		panic(err)
	}

	loaded := &loadedGame{
		openConns: 1,
		game: &Game{
			underlying:     g,
			service:        s,
			uuid:           uuid,
			listenChannels: []chan struct{}{make(chan struct{})},
		},
	}

	loaded.game.mutex.Lock()
	defer loaded.game.mutex.Unlock()

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.games[uuid] = loaded

	updateCh, ok := s.players[playerX]
	if ok {
		go func() {
			updateCh <- NewGameNotification{loaded.game, loaded.game.listenForUpdates()}
		}()
	}

	updateCh, ok = s.players[playerO]
	if ok {
		go func() {
			updateCh <- NewGameNotification{loaded.game, loaded.game.listenForUpdates()}
		}()
	}

	return nil
}

func (s *GameService) OpenGamesForPlayer(playerUUID string) ([]NewGameNotification, <-chan NewGameNotification, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	uuids, err := s.Store.getGameUUIDS(playerUUID)
	if err != nil {
		return nil, nil, err
	}

	count := len(uuids)
	games := make([]NewGameNotification, count)

	for i, uuid := range uuids {
		loaded, ok := s.games[uuid]

		if !ok {
			// load game from db
			underlying, err := s.Store.loadGame(uuid)
			if err != nil {
				return nil, nil, err
			}
			loaded = &loadedGame{
				openConns: 1,
				game: &Game{
					underlying:     underlying,
					service:        s,
					uuid:           uuid,
					listenChannels: []chan struct{}{},
				},
			}

			s.games[uuid] = loaded

			games[i] = NewGameNotification{loaded.game, loaded.game.listenForUpdates()}
		} else {
			// attach to loaded game
			loaded.openConns++
			games[i] = NewGameNotification{loaded.game, loaded.game.listenForUpdates()}
		}
	}

	ch := make(chan NewGameNotification)
	s.players[playerUUID] = ch
	return games, ch, nil
}
