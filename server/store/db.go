package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/google/uuid"
	"github.com/heartles/uttt/server/game"
)

type Player struct {
	UUID, Username, GoogleID string
}

type Store struct {
	db *sql.DB
}

const initUsers = `
CREATE TABLE IF NOT EXISTS "users"
(
    [PK_UUID] TEXT PRIMARY KEY,
    [Username] TEXT UNIQUE NOT NULL,
    [GoogleID] INTEGER UNIQUE NOT NULL
);
`

const initGames = `
CREATE TABLE IF NOT EXISTS "matches"
(
	[PK_UUID] CHAR(36) PRIMARY KEY,
	[GameData] TEXT NOT NULL,
	[UserX] TEXT NOT NULL,
	[UserO] TEXT NOT NULL,
	[Victor] TEXT NOT NULL,
	[LastMoveGameX] INTEGER,
	[LastMoveGameY] INTEGER,
	[LastMoveSubgridX] INTEGER,
	[LastMoveSubgridY] INTEGER,
	[Finished] BOOLEAN NOT NULL,
	FOREIGN KEY (UserX) REFERENCES "users" (PK_UUID),
	FOREIGN KEY (UserO) REFERENCES "users" (PK_UUID),
	FOREIGN KEY (Victor) REFERENCES "users" (PK_UUID)
);
`

func NewStore(filepath string) (*Store, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(initUsers)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(initGames)
	if err != nil {
		return nil, err
	}

	return &Store{db}, nil
}

func (s *Store) TryLookupPlayer(googleID string) (*Player, error) {
	row := s.db.QueryRow(`SELECT PK_UUID, Username FROM users WHERE GoogleID = ?;`, googleID)

	var id string
	var username string
	err := row.Scan(&id, &username)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	fmt.Println("logging in " + id + " " + username)

	return &Player{
		UUID:     id,
		Username: username,
		GoogleID: googleID,
	}, nil
}

func (s *Store) CreatePlayer(username string, googleID string) (*Player, error) {
	id := uuid.New()

	_, err := s.db.Exec(`
		INSERT INTO users (PK_UUID, Username, GoogleID)
		VALUES (?, ?, ?);
		`, id.String(), username, googleID)

	if err != nil {
		// try lookup
		player, _ := s.TryLookupPlayer(googleID)
		if player == nil {
			return nil, err
		}

		return player, nil
	}

	return &Player{
		UUID:     id.String(),
		GoogleID: googleID,
		Username: username,
	}, nil
}

func (s *Store) saveGame(gameID string, game *game.Game) error {
	panic("not impl")
}

func (s *Store) loadGame(gameID string) (*game.Game, error) {
	panic("not impl")
}

func (s *Store) getGameUUIDS(playerID string) ([]string, error) {
	panic("not impl")
}
