package game

import (
	"errors"
	"regexp"
)

// StalematePlayer is the PlayerID for if a stalemate has occurred
const StalematePlayer = "tie"

// ErrSquarePlayed is returned from PlayMove if a given square has
// already been played
var ErrSquarePlayed = errors.New("square already played")

// ErrWrongTurn is returned by PlayMove when a player attempts to play
// a move when it is not their turn
var ErrWrongTurn = errors.New("not this player's turn")

// ErrWrongSubgrid is returned by PlayMove when a player attempts to play
// in the wrong subgrid
var ErrWrongSubgrid = errors.New("incorrect subgrid")

// ErrInvalidPlayer is returned by PlayMove when an invalid player id
// is provided
var ErrInvalidPlayer = errors.New("invalid player id")

// ErrInvalidCoordinate is returned when a coordinate is outside the bounds
// of the game board
var ErrInvalidCoordinate = errors.New("coordinate out of bounds")

// ErrInvalidInput is returned when the game state provided to LoadGame
// is invalid
var ErrInvalidInput = errors.New("invalid game state")

// ErrInvalidLastMove is returned when the last move provided to LoadGame
// is invalid
var ErrInvalidLastMove = errors.New("invalid lastMove")

// SubCoordinate is a reference to either a subgrid or a square
// on the game board
type SubCoordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Coordinate is a reference to a specific square on a specific
// subgrid
type Coordinate struct {
	GameSquare    SubCoordinate `json:"gameSquare"`
	SubgridSquare SubCoordinate `json:"subgridSquare"`
}

// Move encapsulates a player's move
type Move struct {
	PlayerID   string `json:"playerID"`
	Coordinate `json:"coordinate"`
}

// Game encapsulates a game of uttt
type Game struct {
	playerX  string
	playerO  string
	lastTurn *Coordinate
	grid     *subgrid
}

// NewGame is a basic constructor for a Game
func NewGame(playerX, playerO string) (*Game, error) {
	if playerX == playerO {
		return nil, ErrInvalidPlayer
	}

	return &Game{
		playerX, playerO,
		nil,
		initGameBoard(),
	}, nil
}

// LoadGame loads a game from save data
func LoadGame(playerX, playerO string, gameState string, lastTurn *Coordinate) (*Game, error) {
	// remove all invalid characters (useful for whitespace in tests)
	state := regexp.MustCompile("[^XO_]+").ReplaceAllString(gameState, "")
	if len(state) != 3*3*3*3 {
		return nil, ErrInvalidInput
	}

	game, err := NewGame(playerX, playerO)
	if err != nil {
		return nil, err
	}

	err = game.loadState(state, lastTurn)

	return game, err
}

// NewCoordinate is a helper function for constructing a Coordinate
func NewCoordinate(gameX, gameY, subX, subY int) Coordinate {
	return Coordinate{
		SubCoordinate{
			gameX, gameY,
		},
		SubCoordinate{
			subX, subY,
		},
	}
}

func (g *Game) PlayMove(m Move) error {
	player := g.playerIDToEnum(m.PlayerID)
	coord := m.Coordinate
	err := g.verifyMove(player, coord)
	if err != nil {
		return err
	}

	// apply move
	g.grid.board[coord.GameSquare].board[coord.SubgridSquare].state = player
	g.lastTurn = &coord

	// check win condition
	g.grid.board[coord.GameSquare].board[coord.SubgridSquare].match3()
	g.grid.board[coord.GameSquare].match3()
	g.grid.match3()

	return nil
}

func (g *Game) IsCompleted() bool {
	return g.GameWinner() != ""
}

func (g *Game) GameWinner() string {
	return g.playerEnumToID(g.grid.state)
}

func (g *Game) BlockWinner(c SubCoordinate) (string, error) {
	sg, ok := g.grid.board[c]
	if !ok {
		return "", ErrInvalidCoordinate
	}

	return g.playerEnumToID(sg.state), nil
}

func (g *Game) SquareOwner(c Coordinate) (string, error) {
	if !g.isValidCoordinate(c) {
		return "", ErrInvalidCoordinate
	}

	return g.playerEnumToID(g.getSquareState(c)), nil
}

type squareState int

const (
	stateInProgress squareState = iota
	stateTie
	stateX
	stateO
	stateInvalid
)

type game struct {
	playerX  string
	playerO  string
	lastTurn *Coordinate
	grid     *subgrid
}

type subgrid struct {
	// specifies the player that has won (or played in)
	// this subgrid
	state squareState

	// the following is only used if this is a top-level
	// subgrid, they will both have their zero values
	// if this is a bottom-level subgrid
	board map[SubCoordinate]*subgrid
}

func (g *Game) playerEnumToID(p squareState) string {
	switch p {
	case stateInProgress:
		return ""
	case stateO:
		return g.playerO
	case stateX:
		return g.playerX
	case stateTie:
		return StalematePlayer
	default:
		panic("invalid enum value")
	}
}

func (g *Game) playerIDToEnum(p string) squareState {
	if p == g.playerX {
		return stateX
	} else if p == g.playerO {
		return stateO
	}

	return stateInvalid
}

func (g *Game) isValidCoordinate(c Coordinate) bool {
	sg, ok := g.grid.board[c.GameSquare]
	if !ok {
		return false
	}

	_, ok = sg.board[c.SubgridSquare]
	return ok
}

func (g *Game) getSquareState(c Coordinate) squareState {
	return g.grid.board[c.GameSquare].board[c.SubgridSquare].state
}

func (g *Game) loadState(state string, lastTurn *Coordinate) error {
	i := 0
	// GameCoordinate = {{z, w} {x, y}}
	for w := 1; w <= 3; w++ {
		for y := 1; y <= 3; y++ {
			for z := 1; z <= 3; z++ {
				for x := 1; x <= 3; x++ {
					playerChar := state[i]
					i++

					var player squareState
					switch playerChar {
					case '_':
						continue
					case 'X':
						player = stateX
						break
					case 'O':
						player = stateO
						break
					default:
						// invalid chars should've been removed by above regex
						panic("unexpected char found: " + string(playerChar))
					}

					c := NewCoordinate(z, w, x, y)
					g.grid.board[c.GameSquare].board[c.SubgridSquare].state = player
				}
			}
		}
	}
	g.grid.recalcMatch3()

	if lastTurn != nil {
		if !g.isValidCoordinate(*lastTurn) {
			return ErrInvalidCoordinate
		}

		prevPlayer := g.getSquareState(*lastTurn)
		if prevPlayer == stateInProgress {
			return ErrInvalidLastMove
		}

		g.lastTurn = lastTurn
	}

	return nil
}

// verify that a move is valid
func (g *Game) verifyMove(player squareState, coord Coordinate) error {

	if player == stateInvalid {
		return ErrInvalidPlayer
	} else if g.lastTurn == nil {
		// it's the first turn, x goes first
		if player != stateX {
			return ErrWrongTurn
		}
	} else {
		prevPlayer := g.getSquareState(*g.lastTurn)

		if prevPlayer == stateInProgress || prevPlayer == stateInvalid {
			// this should never happen, it would mean the game
			// was in an invalid state. States should be verified
			// at load time
			panic("prevPlayer was none")
		}

		if prevPlayer == player {
			return ErrWrongTurn
		}
	}

	if !g.isValidCoordinate(coord) {
		return ErrInvalidCoordinate
	}

	// you can only play in the subgrid corresponding to the last move
	// UNLESS - it's the first turn OR you are played into a subgrid that's
	// already finished
	if g.lastTurn != nil {
		if g.grid.board[g.lastTurn.SubgridSquare].state == stateInProgress &&
			g.lastTurn.SubgridSquare != coord.GameSquare {
			return ErrWrongSubgrid
		}
	}

	// you can't play in a subgrid that's already won/tied
	if g.grid.board[coord.GameSquare].state != stateInProgress {
		return ErrWrongSubgrid
	}

	// you can't play in a square that's already taken
	if g.getSquareState(coord) != stateInProgress {
		return ErrSquarePlayed
	}

	return nil
}

func (sg *subgrid) at(x, y int) *subgrid {
	return sg.board[SubCoordinate{x, y}]
}

// searches for 3-in-a-row and updates the subgrid status accordingly
func (sg *subgrid) match3() {
	// don't try matching on bottom-level squares
	if sg.board == nil {
		return
	}

	// check horizontals and for open squares
	openSquaresExist := false
	for y := 1; y <= 3; y++ {
		c1 := sg.at(1, y).state
		c2 := sg.at(2, y).state
		c3 := sg.at(3, y).state
		if c1 == stateInProgress || c2 == stateInProgress || c3 == stateInProgress {
			openSquaresExist = true
			continue
		}

		if c1 == c2 && c2 == c3 {
			sg.state = c1
			return
		}
	}

	// check verticals
	for x := 1; x <= 3; x++ {
		c1 := sg.at(x, 1).state
		c2 := sg.at(x, 2).state
		c3 := sg.at(x, 3).state

		if c1 == c2 && c2 == c3 && c1 != stateInProgress {
			sg.state = c1
			return
		}
	}

	// check diagonals
	c1 := sg.at(1, 1).state
	c2 := sg.at(2, 2).state
	c3 := sg.at(3, 3).state

	if c1 == c2 && c2 == c3 && c1 != stateInProgress {
		sg.state = c1
		return
	}

	c1 = sg.at(3, 1).state
	c2 = sg.at(2, 2).state
	c3 = sg.at(1, 3).state

	if c1 == c2 && c2 == c3 && c1 != stateInProgress {
		sg.state = c1
		return
	}

	// neither party has won. the game is either ongoing or stalemate
	// don't modify state if the game is ongoing
	if openSquaresExist {
		return
	}

	// stalemate
	sg.state = stateTie
}

// recursively calls match3 before calling match3 on this grid
func (sg *subgrid) recalcMatch3() {
	if sg.board != nil {
		for _, v := range sg.board {
			v.recalcMatch3()
		}
	}

	sg.match3()
}

func initSubgridBoard() map[SubCoordinate]*subgrid {
	board := map[SubCoordinate]*subgrid{}
	for x := 1; x <= 3; x++ {
		for y := 1; y <= 3; y++ {
			board[SubCoordinate{x, y}] = &subgrid{
				state: stateInProgress,
			}
		}
	}

	return board
}

func initSubgrid() *subgrid {
	return &subgrid{
		state: stateInProgress,
		board: initSubgridBoard(),
	}
}

func initGameBoard() *subgrid {
	grid := map[SubCoordinate]*subgrid{}
	for x := 1; x <= 3; x++ {
		for y := 1; y <= 3; y++ {
			grid[SubCoordinate{x, y}] = initSubgrid()
		}
	}

	return &subgrid{
		stateInProgress,
		grid,
	}
}
