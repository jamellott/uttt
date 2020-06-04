package store

import (
	"database/sql"

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
    [PK_UUID] TEXT UNIQUE PRIMARY KEY,
    [Username] TEXT UNIQUE NOT NULL,
    [GoogleID] INTEGER UNIQUE NOT NULL
);

REPLACE INTO users(PK_UUID, Username, GoogleID) VALUES("tie", "tie", "tie");
`

const initGames = `
CREATE TABLE IF NOT EXISTS "matches"
(
	[PK_UUID] CHAR(36) UNIQUE PRIMARY KEY,
	[GameData] TEXT NOT NULL,
	[UserX] TEXT NOT NULL,
	[UserO] TEXT NOT NULL,
	[Victor] TEXT,
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

	return &Player{
		UUID:     id,
		Username: username,
		GoogleID: googleID,
	}, nil
}

func (s *Store) TryLookupPlayerUUID(id string) (*Player, error) {
	row := s.db.QueryRow(`SELECT GoogleID, Username FROM users WHERE PK_UUID = ?;`, id)

	var googleID string
	var username string
	err := row.Scan(&googleID, &username)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &Player{
		UUID:     id,
		Username: username,
		GoogleID: googleID,
	}, nil
}

func (s *Store) TryLookupPlayerUsername(username string) (*Player, error) {
	row := s.db.QueryRow(`SELECT PK_UUID, GoogleID FROM users WHERE Username = ?;`, username)

	var id string
	var googleID string
	err := row.Scan(&id, &googleID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

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

func (s *Store) saveNewGame(game *game.Game) (string, error) {
	id := uuid.New().String()
	return id, s.saveGame(id, game)
}

func (s *Store) saveGame(gameID string, game *game.Game) error {
	playerX, playerO, state, lastMove := game.SaveGame()

	var lastGameX, lastGameY, lastSubX, lastSubY *int
	if lastMove != nil {
		lastGameX = &lastMove.GameSquare.X
		lastGameY = &lastMove.GameSquare.Y
		lastSubX = &lastMove.SubgridSquare.X
		lastSubY = &lastMove.SubgridSquare.Y
	}

	finished := game.IsCompleted()
	var victor *string
	if finished {
		v := game.GameWinner()
		victor = &v
	}

	_, err := s.db.Exec(`
		REPLACE INTO matches(
			PK_UUID,
			GameData,
			UserX,
			UserO,
			Victor,
			LastMoveGameX,
			LastMoveGameY,
			LastMoveSubgridX,
			LastMoveSubgridY,
			Finished)
			VALUES(?,?.?,?,?,?,?,?,?,?);
		`,
		gameID, state, playerX, playerO, victor, lastGameX, lastGameY,
		lastSubX, lastSubY, finished)
	return err
}

func (s *Store) loadGame(gameID string) (*game.Game, error) {
	row := s.db.QueryRow(`
		SELECT
			GameData,UserX,UserO,
			LastMoveGameX,LastMoveGameY,
			LastMoveSubgridX,LastMoveSubgridY
		FROM matches WHERE PK_UUID = ?;
	`, gameID)

	var state, playerX, playerO string
	var lastGameX, lastGameY, lastSubX, lastSubY *int
	err := row.Scan(&state, &playerX, &playerO, &lastGameX, &lastGameY, &lastSubX, &lastSubY)
	if err != nil {
		return nil, err
	}

	var lastTurn *game.Coordinate
	if lastGameX != nil {
		coord := game.NewCoordinate(*lastGameX, *lastGameY, *lastSubX, *lastSubY)
		lastTurn = &coord
	}

	g, err := game.LoadGame(playerX, playerO, state, lastTurn)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (s *Store) getGameUUIDS(playerID string) ([]string, error) {
	rows, err := s.db.Query(`
		SELECT PK_UUID FROM matches WHERE UserX = ? OR UserO = ?;
	`, playerID, playerID)
	if err != nil {
		return nil, err
	}

	uuids := []string{}
	for rows.Next() {
		uuid := ""
		rows.Scan(&uuid)
		uuids = append(uuids, uuid)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return uuids, nil
}
