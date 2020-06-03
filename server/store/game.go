package store

import (
	"errors"
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

func (g *Game) Close(ch <-chan struct{}) error {
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
		return errors.New("Invalid channel")
	}

	close(g.listenChannels[idx])
	g.listenChannels = append(g.listenChannels[:idx], g.listenChannels[idx:]...)

	g.mutex.Unlock()
	return g.service.closeGame(g)
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
	games map[string]*loadedGame
	mutex sync.Mutex
	*Store
}

func NewGameService(dbFilename string) (*GameService, error) {
	st, err := NewStore(dbFilename)
	if err != nil {
		return nil, err
	}

	return &GameService{
		map[string]*loadedGame{},
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

func (s *GameService) OpenGamesByPlayer(playerUUID string) ([]*Game, []<-chan struct{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	uuids, err := s.Store.getGameUUIDS(playerUUID)
	if err != nil {
		return nil, nil, err
	}

	count := len(uuids)
	games := make([]*Game, count)
	chans := make([]<-chan struct{}, count)

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

			games[i] = loaded.game
			chans[i] = loaded.game.listenForUpdates()
		} else {
			// attach to loaded game
			loaded.openConns++
			games[i] = loaded.game
			chans[i] = loaded.game.listenForUpdates()
		}
	}
	return games, chans, nil
}
