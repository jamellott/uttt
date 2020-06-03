package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/heartles/uttt/server/game"
)

type Store struct {
	db *sql.DB
}

const initCmd = `
CREATE TABLE IF NOT EXISTS "users"
(
    [PK_UUID] CHAR(36) PRIMARY KEY NOT NULL,
    [Name] NVARCHAR(32) NOT NULL,
    [GoogleID] INTEGER NOT NULL,
);

CREATE TABLE IF NOT EXISTS "matches"
(
	[PK_UUID] CHAR(36) PRIMARY KEY NOT NULL,
	[GameData] TEXT NOT NULL,
	FOREIGN KEY (UserX) REFERENCES users(PK_UUID) NOT NULL,
	FOREIGN KEY (UserO) REFERENCES users(PK_UUID) NOT NULL,
	[LastMoveGameX] INT,
	[LastMoveGameY] INT,
	[LastMoveSubgridX] INT,
	[LastMoveSubgridY] INT,
	[Finished] BOOL NOT NULL,
	FOREIGN KEY (Victor) REFERENCES users(PK_UUID),
);
`

func NewStore(filepath string) (*Store, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	return &Store{db}, nil
}

func (s *Store) LookupPlayerUUID(username string) (string, error) {
	row := s.db.QueryRow(`SELECT PK_UUID FROM users WHERE Name = ?;`)

	var uuid string
	err := row.Scan(&uuid)
	if err == sql.ErrNoRows {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return uuid, nil
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
