package store

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

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

	err := g.underlying.PlayMove(m)
	for _, ch := range g.listenChannels {
		go func(ch chan struct{}) {
			g.mutex.RLock()
			defer g.mutex.RUnlock()
			ch <- struct{}{}
		}(ch)
	}
	return err
}

// The write mutex must be held during this call
func (g *Game) listenForUpdates() <-chan struct{} {
	ch := make(chan struct{})
	g.listenChannels = append(g.listenChannels, ch)

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
		notif := NewGameNotification{loaded.game, loaded.game.listenForUpdates()}
		go func(ch chan NewGameNotification, notif NewGameNotification) {
			ch <- notif
		}(updateCh, notif)
	}

	updateCh, ok = s.players[playerO]
	if ok {
		notif := NewGameNotification{loaded.game, loaded.game.listenForUpdates()}
		go func(ch chan NewGameNotification, notif NewGameNotification) {
			ch <- notif
		}(updateCh, notif)
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
			// don't need the game mutex held here because nobody else can have
			// a handle to it yet
			games[i] = NewGameNotification{loaded.game, loaded.game.listenForUpdates()}
			go s.periodicFlushToDB(loaded.game)
		} else {
			// attach to loaded game
			loaded.openConns++
			loaded.game.mutex.Lock()
			games[i] = NewGameNotification{loaded.game, loaded.game.listenForUpdates()}
			loaded.game.mutex.Unlock()
		}
	}

	ch := make(chan NewGameNotification)
	s.players[playerUUID] = ch
	return games, ch, nil
}

func (s *GameService) periodicFlushToDB(g *Game) {
	for {
		<-time.After(1 * time.Minute)
		g.mutex.RLock()
		if len(g.listenChannels) == 0 {
			// game has been unloaded, do nothing
			// and exit
			return
		}

		err := s.Store.saveGame(g.uuid, g.underlying)
		if err != nil {
			fmt.Println(err)
		}

		g.mutex.RUnlock()

	}
}
